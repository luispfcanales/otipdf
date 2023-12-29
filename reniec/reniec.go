package reniec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
	args["protocol"] = "T"                                                          //https: S - http: T
	args["fileDownloadUrl"] = "http://18.219.214.89:4000/reniec/download"           //endpoint
	args["fileDownloadLogoUrl"] = ""                                                //logo
	args["fileDownloadStampUrl"] = "http://18.219.214.89:4000/public/logofirma.png" //stamp reniec logo - optional
	args["fileUploadUrl"] = "http://18.219.214.89:4000/file/upload"                 //route to upload file and save
	args["contentFile"] = opts.FileID + ".pdf"                                      //real name document - json struct
	args["reason"] = opts.Reason                                                    //json struct
	args["pageNumber"] = opts.PageNumber                                            //json struct
	//args["posx"] = "339.5"                                                      //json sctruct
	//args["posy"] = "658.2"                                                      //json sctruct
	args["posx"] = opts.Pox //json sctruct
	args["posy"] = opts.Poy //json sctruct
	args["isSignatureVisible"] = "true"
	args["stampAppearanceId"] = opts.StampAppearanceID //json struct
	args["fontSize"] = "7"
	args["dcfilter"] = ".*FIR.*|.*FAU.*"
	//args["signatureLevel"] = "0" //pcx why info set 0?
	args["outputFile"] = opts.FileID + "[R].pdf" //json struct name file
	args["maxFileSize"] = "41943040"             //40Mb
	args["timestamp"] = "false"
	log.Println(topoint(opts.Pox))
	log.Println(topoint(opts.Poy))
	rs, err := json.Marshal(args)
	if err != nil {
		log.Println(err)
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(rs)
	log.Println("se entrega el credenciales en base64")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"args": encodedData,
	})
}

func topoint(pos string) string {
	value, err := strconv.ParseFloat(pos, 64)
	if err != nil {
		return "0"
	}
	value = value / 0.352777778
	return fmt.Sprintf("%v", value)
}

func LoadFirm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("se inicia proceso de subida de documento")
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
	log.Println("archivo:", f)
	log.Println(h.Filename) //[outputFile]
	fl, err := os.Create("/mnt/s3/ServicesSheet/" + h.Filename)
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
	log.Println("se subio el documento firmado->", b)

	w.WriteHeader(200)
	w.Write([]byte(""))

}

func DownloadFirm(w http.ResponseWriter, r *http.Request) {
	log.Println("inicia proceso de descarga del documento")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fn := "testhoja.pdf"
	fs, err := os.Open("/mnt/s3/ServicesSheet/" + fn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer fs.Close()

	w.Header().Add("Content-Type", "application/pdf")
	w.Header().Add("Content-Type", "filename="+fn)

	log.Println("termina proceso de descarga del documento")
	io.Copy(w, fs)
	w.WriteHeader(200)
}

func DownloadReniec(w http.ResponseWriter, r *http.Request) {

	log.Println("inicia proceso de descarga del documento")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fn := "testhoja.pdf"
	fs, err := os.Open("/mnt/s3/ServicesSheet/" + fn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer fs.Close()

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename="+fn)

	log.Println("termina proceso de descarga del documento")
	io.Copy(w, fs)
	w.WriteHeader(200)
}
