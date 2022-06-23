package url_args

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_compareLinkWithArgs(t *testing.T) {
	t.Run("not has url args", func(t *testing.T) {
		assert.True(t, Compare(
			"/abcd/sbxwws?abcd=100&x=111",
			"/abcd/sbxwws?abcd=100&x=112",
			[]string{"abcd"},
		))
		assert.False(t, Compare(
			"/abcd/sbxwws?abcd=100&x=111",
			"/abcd/sbxwws?abcd=100&x=112",
			[]string{"abcd", "x"},
		))
	})

	t.Run("has url args and link eq", func(t *testing.T) {
		assert.True(t, Compare(
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd",
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd",
			[]string{"url.sn", "url.__biz"},
		))
	})

	t.Run("has url args and not eq", func(t *testing.T) {
		assert.False(t, Compare(
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzwkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd",
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd",
			[]string{"url.sn", "url.__biz"},
		))
	})

	t.Run("has url args and args eq", func(t *testing.T) {
		assert.True(t, Compare(
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F_biz%3Dxx%26__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad",
			"/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd",
			[]string{"url.sn", "url.__biz"},
		))
	})
}
