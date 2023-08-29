package reniec

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type OptionsSignature struct {
	StampAppearanceID string `json:"stamp_appearance_id,omitempty"`
	FileID            string `json:"file_id,omitempty"`
	PageNumber        string `json:"page_number,omitempty"`
	Pox               string `json:"pox,omitempty"`
	Poy               string `json:"poy,omitempty"`
	Reason            string `json:"reason,omitempty"`
}

func GetArgs(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		return
	}

	opts := &OptionsSignature{}
	err := json.NewDecoder(r.Body).Decode(opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	args := make(map[string]string)

	args["app"] = "pdf" //pcx
	//args["mode"] = "lot-p"//pcx
	args["clientId"] = "ZIzAvpCQernywPNktelaHQH0yi0"
	args["clientSecret"] = "B6jWcQmOjJkD94A-EgTl"
	args["idFile"] = "load_file"
	args["type"] = "W"
	args["protocol"] = "T"                                                        //https: S - http: T
	args["fileDownloadUrl"] = "http://18.118.181.184/file/download"               //endpoint
	args["fileDownloadLogoUrl"] = ""                                              //logo
	args["fileDownloadStampUrl"] = "http://18.118.181.184/public/logo-unamad.png" //stamp reniec logo - optional
	args["fileUploadUrl"] = "http://18.118.181.184/file/upload"                   //route to upload file and save
	args["contentFile"] = opts.FileID                                             //real name document - json struct
	args["reason"] = opts.Reason                                                  //json struct
	args["pageNumber"] = opts.PageNumber                                          //json struct
	args["posx"] = opts.Pox                                                       //json sctruct
	args["posx"] = opts.Poy                                                       //json sctruct
	args["isSignatureVisible"] = "true"
	args["stampAppearanceId"] = opts.StampAppearanceID //json struct
	args["fontSize"] = "7"
	args["dcfilter"] = ".*FIR.*|.*FAU.*"
	//args["signatureLevel"] = "0" //pcx why info set 0?
	args["outputFile"] = "firmado.pdf" //json struct name file
	args["maxFileSize"] = "41943040"   //40Mb
	args["timestamp"] = "false"

	rs, err := json.Marshal(args)
	if err != nil {
		log.Println(err)
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(rs)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"args": encodedData,
	})
}

func LoadFirm(w http.ResponseWriter, r *http.Request) {
	log.Println("se ejecuta")
	err := r.ParseMultipartForm(40 << 20)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	f, h, err := r.FormFile("load_file")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//url.QueryUnescape
	fl, err := os.Create(h.Filename + ".pdf")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	defer fl.Close()

	b, err := io.Copy(fl, f)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	log.Println("write bytes are:", b)

	w.WriteHeader(200)
	w.Write([]byte(""))

}

func DownloadFirm(w http.ResponseWriter, r *http.Request) {
	log.Println("entro download")
	fn := "38be5475-6b48-4dd9-83fd-77f51dfdb97e.pdf"
	fs, err := os.Open("tmp/" + fn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer fs.Close()

	w.Header().Add("Content-Type", "application/pdf")
	w.Header().Add("Content-Type", "filename="+fn)

	io.Copy(w, fs)
	w.WriteHeader(200)
}
