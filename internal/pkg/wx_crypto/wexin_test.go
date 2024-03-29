package wx_crypto

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAesCBCEncrypt(t *testing.T) {
	sessionKey := "z81Q1p9jPDumiyz2Vi33zQ=="
	encryptedData := "qbfszwvVbnsWYkZcvLYMQWHVvo1lPpy+0Z3CzipxTUl/hN2MEbVghUaWXrDS0N3Edjxvhhi27uOBAtMMhnZqp4QHrRxTyAXuyk6bB7aB01lWbNTvqqffmnTtaAzk5BwVaPFiyRP2aoOmGaKQGjlVxCHU7UJ0N1dAUAcSTKB3f4fgfOy79TYK41AOB+9GFDaW85mA6wsnmek5jQ05NLNTJg=="
	iv := "raUOWdl0H3/ORd9wSbKrRQ=="

	//appID := "wx4053ee583576b3c7"
	//sessionKey := "cSeKHbGPd1zGCeftm0cveg=="
	//encryptedData := "FJtT4sSAE/jjlrkAs5ZqZEfwaes1+/zJiIyQfY5kZ+sw4qtrt+VTz+JUNgJsKKA19HVPMAGyeRD+HIT4cTDXOPmXfc5HOOx4L87fOJeG83d8jQ7+iWVBe8G6ZifcA4Wo5uaq4AihgfiBPIxxc5FycPpzx2DHgSyOxl7KGxZ0omzgtrfIdRftkRuKkdL/X+hKnPKECiSQvvAwmgLG9XWzFP29X/o4gtrsDex0Q6UFNjogmBn28E/6TIbmXDK494Z791dzk3m4WvN6akX8va2GsN4JTDEaOZTnO/SX/PfLzSLPrC2KmGbEQDLeOdoBNhDu+X1k68VHa3wImc+KLISKRHopx2tX56yAxM5vL6i1S3h68GtfToNdRqHV1g6x+akGSnA4F0HactXdjq5+A3hTgIovi3iAvhOTN6yQjVT4iWk2YCPsJ42VZwuAssVrsP7Y16nDyoTWhCuygKLqib84M5nmmJNoMLIbmuB9vW7bQ9DqAq0lnZ/ZKwXDtfas89XUoZLRulzsiDxtSTr7qYCWu4VkINOROT7TV+MiamntrhFh7L9xRwOt/SYK9fZlVnpvpwOjGGWa6YvFiq5R5ScGkP4L12yEHGQzI4KLGcXmXK+PnOb0i+sgErXyrlVXB20me65BKfrjtDwtpDKNgsbpE2Pdi+nmcbhrUpop0YgJjFIHCCW3x+ZzkxZX4qKzpyzHQDEOp8EVUEI8vb6FRGaZF5S+k/mOEknjGZssOTiX+BsuBJ9r5XsHHQRa539tX9EurxJAmcbou+S6bjDx+94kLivx3+xvOaJ+56U8sdx31b40F87Wrd//zaRyELAOsho7YShFTfmb0vyprNUJtTTC8s6MuxnZco+GwiSd+iDfYY7vmJhFY3pgBu48RLYDACVjiHucDWUMh2Y6DWSNEKDPtUWy/mWrzpu9W8YUjV+q8EfiPG3P/A9T6FTlwOBe0/L0rbdZGjoqhYwOSqKS3RJvkLad535aG5A8gcH0YXH42Ze6OO+LkKUOOpz/XaWThQj+voxXz2Pbi0tZXRG+mKvBeftP9PrhJpN4AmoWCPd5CA6yJC+yT50F02XpyXJXqQVGpuz2CPkKGNh8ASQtldldL2+eiHGyl9m2S3vgk6hNCHl48/Dx9txoywkCY4DWIbx5SIA3LgNdqT9g1AwGPtlI9o0R6asbLUFG3dC0jSARAJTGp61KrLoWXxACp2hZ9Heha8uLsjPK5a8v4h06FP6VfU0epN3VsHC2/nfgT7TXaY03BSI7hZuNRBXUpId9dbVG5fueEXl3p5TrupM3k6QNr+hny90r/F//s5UeQAnhkYeE8pm5/m3HcLNKc19UysdBr3zQN0vW9ro8+gi5cUR0FtTOReWCCGEECebWBKh1a6NVfuQFrv8E0jFcwrkdiRch0H9vvSM4JnNWAU98K7vyR7yRblLN8MZp5mAHJdqzGn+OiPiRaK4pRKk2ZtHspYaBT2WT6RgVEDWNbmLa8n2gJCspM+NS+mjcGfqIDr19M/LEIHFRyLnoMTS+25tOqbIkB6lEsMyj3U5lLjrYMmaf9qzrZw+KIfduohBLciJEcuXSXPBlchPJ/A4l/2OPeNOGUcrC9QI1Jcd1nIRPsANE5+LYS5bMLNa0KLvIjiAvUGXonFYyanTs30xtXGguqtoR"
	//iv := "NuM7yNcjmjdtFmjxHk4g5A=="

	var data WxUserPhone
	s, err := Decrypt(encryptedData, iv, sessionKey, &data)
	assert.Nil(t, err)
	fmt.Printf("%s\n", s)

	iv2, err := base64.StdEncoding.DecodeString(iv)
	require.Nil(t, err)
	encrypt, err := Encrypt(sessionKey, string(iv2), s)
	assert.Nil(t, err)
	fmt.Printf("%s\n", encrypt)

	assert.Equal(t, encryptedData, encrypt)
}

func TestPkcs7(t *testing.T) {
	d := "z81Q1p9jPDumiyz2Vi33zQ=="
	pad, err := pkcs7Pad([]byte(d), 7)
	require.Nil(t, err)
	unpad, err := pkcs7Unpad(pad, 7)
	require.Nil(t, err)
	assert.Equal(t, d, string(unpad))
}
