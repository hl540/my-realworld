package util

import "testing"

const testSecretKey = "your-256-bit-secret"

func TestNewJwtByToken(t *testing.T) {
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYmIiOjExMTExMTExMTExMTEsInh4eCI6Inh4eHh4eHh4eHh4eCJ9.lZzzQNKq8HzONmzW8k5zzg4DFERnehQ1zbsfQI8mBlU"
	jwt, err := NewJwtByToken(testSecretKey, str)
	if err != nil {
		t.Error(err)
		return
	}
	a := jwt.GetString("xxx")
	b := jwt.GetInt("bbb")
	t.Log(a, b)
	s, err := jwt.Token()
	t.Log(s, err)
}

func TestNewJwtByData(t *testing.T) {
	jwt := NewJwtByData(testSecretKey, map[string]interface{}{
		"xxx": "xxxxxxxxxxxx",
		"bbb": 1111111111111,
	})
	a := jwt.GetString("xxx")
	b := jwt.GetInt("bbb")
	t.Log(a, b)
	s, err := jwt.Token()
	t.Log(s, err)
}
