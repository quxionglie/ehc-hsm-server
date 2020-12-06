package utils

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func TestMindex(t *testing.T) {
	log.Printf("%X", Sm3("123"))
}

// 健康卡计算国家规范
func TestEhcId(t *testing.T) {
	//        电子健康卡管理信息系统密码服务功能区域包含应用密码机，为电子健康卡管理信息系 统的二维码管理模块提供密码服务。
	//        应用密码机是电子健康卡管理信息系统的一部分，在物 理环境上应与电子健康卡管理信息系统的其他功能模块同区域部署;
	//        应用加密机仅对本电子 健康卡管理信息系统的二维码管理模块提供密码服务，不向其他任何系统提供服务。
	//        应用密码机提供统一的主索引ID生成接口和电子健康卡ID接口，电子健康卡管理信息系 统调用应用密码机模块时，
	//        应符合附录D中的“密码机高级应用编程接口”的要求。
	//        应用密码机中包含安全密钥和保护密钥。
	//        1) 安全密钥
	//        安全密钥用于生成电子健康卡ID。 电子健康卡管理信息系统生成用户主索引ID后，调用密码服务功能生成电子健康卡ID。
	//        以应用加密机中的安全密钥为根密钥;将主索引ID以8字节为一组，将各组异或后的结果作 为分散因子;
	//        将根密钥通过SM4算法分散产生用户身份认证密钥，分散方案与实体卡密钥分 散方案一致。
	//        通过用户身份认证密钥和SM4算法对用户证件类型和证件号码加密生成电子健 康卡ID。电子健康卡ID采用以下表达式生成:
	//        电子健康卡ID=SM4加密(证件类型+证件号，用户身份认证密钥，ECB模式)
	//        2) 保护密钥
	//        保护密钥用于动态二维码生成与验证。生成动态二维码时，二维码管理调用密码服务， 通过SM4算法将有效时间的明文信息加密。验证动态二维码时，
	//        二维码管理调用密码服务， 通过SM4算法将有效时间的密文解密，进行验证。
	//        有效时间采用16进制编码，加密补位为0。格式为YYYYMMDDHHMMSS。如当前时间为2018 年1月17日11:50:03，则表示为20180117115003(HEX)。
	//        二维码有效性信息=SM4加密(有效时间补位结果，保护密钥，ECB模式)
	//java运行结果20201207：
	//分散因子(计算)=92CDCE45A16799FD
	//分散因子(规范)=92CDCE45A16799FD
	//分散因子(取反)=6D3231BA5E986602
	//用户身份认证密钥(计算)=86C63180C2806ED1F47B859DE501215B
	//用户身份认证密钥(规范)=86C63180C2806ED1F47B859DE501215B
	//idNoBytes=3031313130313030323031373132323530303658000000000000000000000000
	//健康卡(计算1)=98FF9F2C05145CB9305E9D8A57E072BB51180CA49E5BE5E9D41BCE3A9A571928
	//健康卡(规范1)=98FF9F2C05145CB9305E9D8A57E072BB51180CA49E5BE5E9D41BCE3A9A571928
	//健康卡(检验值)=3031313130313030323031373132323530303658000000000000000000000000
	//健康卡(规范值)=3031313130313030323031373132323530303658000000000000000000000000

	var mkey []byte
	mindex := "ABD17E7ED399EF68AB5660155D6E226D2C92EAC3254A4A66BED83AED0ADA1E9E"
	for i := 0; i < 64; i = i + 16 {
		tmp := mindex[i : i+16]
		tmpByte, _ := hex.DecodeString(tmp)
		if mkey == nil {
			mkey = tmpByte
		} else {
			var x = BytesToInt64(mkey) ^ BytesToInt64(tmpByte)
			mkey = Int64ToBytes(x)
		}
	}

	log.Println("分散因子(计算)=" + strings.ToUpper(hex.EncodeToString(mkey)))

	fx := "92CDCE45A16799FD"
	var fxByte, _ = hex.DecodeString(fx)
	fxBytef := make([]byte, len(fxByte))
	var temp byte
	for i := 0; i < len(fxByte); i++ {
		temp = fxByte[i]
		fxBytef[i] = (byte)(^temp)
	}

	// 6D3231BA5E986602
	x := strings.ToUpper(hex.EncodeToString(fxBytef))
	log.Println("分散因子(规范)=" + fx)
	log.Println("分散因子(取反)=" + x)

	//92CDCE45A16799FD6D3231BA5E986602
	//byte[] concatBytes = ArrayUtils.addAll(fxByte, fxBytef);
	concatBytes := append(fxByte, fxBytef...)

	key := "0123456789ABCDEFFEDCBA9876543210"
	safeKey, _ := hex.DecodeString(key)

	// byte[] x1 = SM4.encrypt_Ecb_nopadding(safeKey, concatBytes);
	x1, err3 := Sm4EncryptNoPadding(safeKey, concatBytes)
	if err3 != nil {
		t.Fail()
	}
	//  System.out.println(SM4.encryptEcb(key, Hex.encodeHexString(concatBytes)));
	log.Println("用户身份认证密钥(计算)=" + strings.ToUpper(hex.EncodeToString(x1)))
	log.Println("用户身份认证密钥(规范)=" + "86C63180C2806ED1F47B859DE501215B")

	idNoStr := "01" + "11010020171225006X"
	idNoBytesOri := []byte(idNoStr)
	idNoBytes := make([]byte, 32)
	for i := 0; i < len(idNoBytes); i++ {
		//idNoBytes[i] = 0;
		if i < len(idNoBytesOri) {
			idNoBytes[i] = idNoBytesOri[i]
		}
	}

	//idNoBytes=3031313130313030323031373132323530303658000000000000000000000000
	log.Println("idNoBytes=" + strings.ToUpper(hex.EncodeToString(idNoBytes)))
	keyFactors, _ := hex.DecodeString("86C63180C2806ED1F47B859DE501215B")
	//如果是nopadding则传入的字节数据长度应该是16的倍数
	idNoBytes = ZeroPadding(idNoBytesOri, 16)
	ehcIdEncBytes, _ := Sm4EncryptNoPadding(keyFactors, idNoBytes)

	ehcId := strings.ToUpper(hex.EncodeToString(ehcIdEncBytes))
	ehcIdStd := "98FF9F2C05145CB9305E9D8A57E072BB51180CA49E5BE5E9D41BCE3A9A571928"
	log.Println("健康卡(计算1)=" + ehcId)
	log.Println("健康卡(规范1)=" + ehcIdStd)
	assert.Equal(t, ehcIdStd, ehcId)

	// 86C63180C2806ED1F47B859DE501215B
	// 86C63180C2806ED1F47B859DE501215B002A8A4EFA863CCAD024AC0300BB40D2

	decodeIdNoBytes, _ := Sm4DecryptNoPadding(keyFactors, ehcIdEncBytes)
	log.Println("健康卡(检验值)=" + strings.ToUpper(hex.EncodeToString(decodeIdNoBytes)))
	log.Println("健康卡(规范值)=" + "3031313130313030323031373132323530303658000000000000000000000000")

}

func TestEncTime(t *testing.T) {
}

func TestDecTime(t *testing.T) {
}
