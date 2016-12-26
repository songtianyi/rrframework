package storage

type Ufile struct {
	PublicKey string
	PrivateKey string
	BucketName string
}

const (
	EXPIRE = 3600
	SUFFIX = ".ufile.ucloud.cn"
)

func (s *Ufile) signature(method string) {
	expir := strconv.FormatInt(EXPIRE + time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond))/1000, 10)
	data := method + "\n"
	data += "\n" //Content-md5 null
	data += "\n" //Content-Type null
	data +=  expir + "\n"
	data += "/" + s.BucketName + "/" + s.PublicKey

	h := hmac.New(sha1.New, []byte(PRIVATE_KEY))
	h.Write([]byte(data))
	sEnc := b64.StdEncoding.EncodeToString(h.Sum(nil))
	return sEnc	
}

func (s *Ufile) Save() {
	signature("POST")
}

func (s *Ufile) Fetch() [
}
