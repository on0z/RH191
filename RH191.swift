import Foundation

/**
 三菱エアコンのリモコンのエミュレータ
 */

/// オプション定義

/** active
 0: off
 1: on
*/
var active: Int = 0

/** mode
 TargetHeaterCoolerState
 0: Auto
 1: Heat
 2: Cool
 */
var mode: Int = 0

/** 温度
 16~31
 */
var temperature: Int = 20

/// コマンドライン引数処理

let argv: [String] = ProcessInfo.processInfo.arguments
guard argv.count == 4 else {
    print("failed")
    exit(-1)
}

let actives: [String] = [
    "0", "INACTIVE",
    "1", "ACTIVE"
]
active = Int((actives.firstIndex(of: argv[1].uppercased()) ?? -2) / 2)
guard 0 <= active && active <= 1 else {
    print("failed")
    exit(-1)
}

let modes: [String] = [
    "0", "AUTO",
    "1", "HEAT",
    "2", "COOL"
]
mode = Int((modes.firstIndex(of: argv[2].uppercased()) ?? -2) / 2)
guard 0 <= mode && mode <= 3 else {
    print("failed")
    exit(-1)
}

temperature = Int(argv[3]) ?? -1
guard 16 <= temperature && temperature <= 31 else {
    print("failed")
    exit(-1)
}

/// ジェネレータ

extension String{
    func split(by length: Int) -> [String]{
        var tmpString = ""
        var splitted: [String] = []
        for (i, c) in self.enumerated(){
            tmpString.append(c)
            if (i + 1)%length == 0{
                splitted.append(tmpString)
                tmpString = ""
            }
        }
        if tmpString != ""{
            splitted.append(tmpString)
        }
        return splitted
    }
}


enum Active: Int { case off, on }
enum Mode: Int { case dry, heat, cool }
enum Speed: Int { case auto, weak, middle, strong }
enum Direction: Int { case auto, up, middleUp, middle, middleDown, down, unknown, move }

/// 送信データのHexを生成する
/// - Parameters:
///   - active: 電源
///   - mode: 運転モード
///   - temperature: 温度
///   - speed: 風速
///   - direction: 風向
/// - Returns: 送信データのHex
func genHex(active: Active, mode: Mode, temperature: Int, speed: Speed, direction: Direction) -> String{
    var bytes: [String] = ["23", "CB", "26", "01", "00"]
    switch active {
    case .on:
        bytes.append("20")
    case .off:
        bytes.append("00")
    }
    
    /// 温度guard
    guard 16 <= temperature && temperature <= 31 else{
        return ""
    }
    
    /// モードコード
    switch mode{
    case .dry:
        bytes.append("10")
    case .heat:
        bytes.append("08")
    case .cool:
        bytes.append("18")
    }
    
    /// 温度
    if mode == .dry{
        bytes.append("08")
    }else{
        bytes.append(
            String(format: "%02X", temperature - 16)
        )
    }
    
    /// モードコード
    switch mode{
    case .dry:
        bytes.append("32")
    case .heat:
        bytes.append("30")
    case .cool:
        bytes.append("36")
    }
    
    var speedDirByte = 0
    speedDirByte += speed.rawValue
    speedDirByte += direction.rawValue << 3
    speedDirByte += 0b01 << 6
    bytes.append(String(format: "%02X", speedDirByte))
    
    bytes += ["00", "00", "00", "00", "00", "00", "00"]
    
    var sum = 0
    for byte in bytes{
        sum += Int(byte, radix: 16) ?? 0
    }
    sum &= 0xFF
    
    let code: String = bytes.joined() + String(format: "%2X", sum)
    
    return code
}

/// 各バイトをバイナリベースで逆順にする
/// - Parameter hex: 対象の16進数文字列
/// - Returns: 逆順にされた文字列
func revHex(hex: String) -> String{
    var revStr: String = ""
    for byte in hex.split(by: 2){
        /// まず数値化する number
        guard let n: Int = Int(byte, radix: 16) else {
            return ""
        }
        
        /// 数値を二進数の文字列にする binary
        var b = String(n, radix: 2)
        b = [String](repeating: "0", count: (8 - b.count)).joined() + b
        
        /// 逆順にする
        b = String(b.reversed())
        
        /// bを数値にする reversed number
        guard let rn: Int = Int(b, radix: 2) else {
            return ""
        }
        
        /// 数値を16進数の文字列にする reversed hex
        let rh = String(format: "%02X", rn)
        
        revStr += rh
    }
    return revStr
}


/// HexからADRSIRコードを生成する
/// - Parameter hex: 送信するHex
/// - Returns: ADRSIRコード
func genADRSIR(hex: String) -> String{
    /// 始端
    var adrsir: String = "83004300"
    
    /// Hexをバイト単位に分割する
    for byte in hex.split(by: 2){
        /// バイトを二進数表記にする %08b
        var binary = String(Int(byte, radix: 16) ?? 0, radix: 2)
        binary = [String](repeating: "0", count: (8 - binary.count)).joined() + binary
        
        /// 1と0をADRSIR表記にする
        for bit in binary{
            if bit == "1" {
                adrsir += "11003200"
            }else if bit == "0"{
                adrsir += "11001100"
            }
        }
    }
    /// 終端
    adrsir += "1300EE01"
    
    /// 2回繰り返して返却する
    return adrsir + adrsir
}

/// テスト用
enum RemoconFormat{case AEHA, NEC, Sony, Unknown}
func convertAdrsirToHex(ADRSIR: String, format: RemoconFormat = .AEHA) -> String{
    /// 4文字ごとに区切る
    let pulseCounts: [String] = ADRSIR.split(by: 4)
    guard pulseCounts.count % 2 == 0 else {
        print("pulseCounts.countがおかしい")
        exit(2)
    }

    /// 元データは38kHzパルスの数を2バイトずつのリトルエンディアンで表しているので，入れ替える必要がある．
    /// "8400" → 0084H → 132d
    /// 賢い人なら，前から2文字ずつ区切って先頭から直に10進数化するかもね．
    var formattedPulseCounts: [Int] = []
    for pulseCount in pulseCounts{
        /// 簡単に説明すると，pulseCountは絶対に4文字のはずです．
        /// それを2文字と2文字に分割し，逆にして結合した文字列を作ります．
        /// それをInt(_: radix:)に放り込んで10進数にし，配列に放り込んでいます．
        formattedPulseCounts.append(
            Int(
                pulseCount.split(by: 2).reversed().joined(),
                radix: 16
            ) ?? 0
        )
    }
    
    if format == .Unknown{
        print(formattedPulseCounts)
    }

    if format == .AEHA || format == .NEC{
        /// パルス数をビットに変換していく
        /// ONが1T，OFFが1T続けば0．ONが1T，OFFが3T続けば1となる
        /// NECフォーマットでは、ONが16T、OFFが8T続いた後、上記のデータ部が始まります。一方、家製協(AEHA)フォーマットでは、ONが8T、OFFが4T続いた後、上記のデータ部が始まります。
        var bits: String = ""
        for i in stride(from: 0, to: formattedPulseCounts.count, by: 2){
            // print(formattedPulseCounts[i], formattedPulseCounts[i + 1])
            if formattedPulseCounts[i] > 100 {
                /// leader
                /// 本当はここでTの値を求めるべき
                bits.append("11111111")
                continue
            }
            if formattedPulseCounts[i+1] > 100 {
                /// leader
                bits.append("11111111")
                /// 本当はここでTの値を求めるべき
                continue
            }
            if (formattedPulseCounts[i + 1] / formattedPulseCounts[i]) >= 2{
                /// オフ(i+1)がオンの3倍なら1
                bits.append("1")
            }else{
                bits.append("0")
            }
        }

        var bytes: String = ""
        for byte in bits.split(by: 8){
            // print(byte)
            bytes.append(String(format: "%02X", Int(byte, radix: 2) ?? 0))
        }

        return bytes
    }
    return ""
}

/// Hexを生成する
let gennedhex = genHex(
    active: Active(rawValue: active) ?? .off,
    mode: Mode(rawValue: mode) ?? .dry,
    temperature: temperature,
    speed: .auto,
    direction: .auto)
/// 反転させる
let revedHex = revHex(hex: gennedhex)
//print(revedHex)
/// ADRSIR化する
let gennedAdrsir = genADRSIR(hex: revedHex)
//print(convertAdrsirToHex(ADRSIR: gennedAdrsir))
/// 出力
print(gennedAdrsir)
