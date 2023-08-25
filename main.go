package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/luispfcanales/otipdf/models"
	"github.com/luispfcanales/otipdf/reniec"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

const TITLE_PAGE string = "Hoja de servicio"

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/args", reniec.GetArgs)
	http.HandleFunc("/pdf", response)
	http.HandleFunc("/firm", vista)
	http.HandleFunc("/file/upload", reniec.LoadFirm)
	http.HandleFunc("/file/download", reniec.DownloadFirm)

	http.HandleFunc("/ext", info)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	http.ListenAndServe(":"+port, nil)
}
func vista(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("vista.html")
	if err != nil {
		log.Println(err)
		return
	}
	t.Execute(w, nil)
}

func response(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	return
	//}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetTitle(TITLE_PAGE, true)

	headerTable := []string{"", "UNIVERSIDAD NACIONAL AMAZONICA DE MADRE DE DIOS", ""}
	contentTable := [][]string{
		{"", `"Madre de Dios Capital de la Biodiversidad del Peru"`, ""},
		{"", `"Ano de la Unidad, la Paz y el Desarrollo"`, ""},
	}

	m.Row(16, func() {
		m.Col(1, func() {
			_ = m.FileImage("logo-unamad/logo-unamad.png", props.Rect{
				Left:    0,
				Top:     0,
				Center:  true,
				Percent: 100,
			})
		})
		m.ColSpace(10)
		m.Col(1, func() {
			_ = m.FileImage("logo-unamad/logo-unamad.png", props.Rect{
				Left:    0,
				Top:     0,
				Center:  true,
				Percent: 100,
			})
		})
		m.TableList(headerTable, contentTable, props.TableList{
			HeaderProp: props.TableListContent{
				Family:    consts.Arial,
				Style:     consts.Bold,
				Size:      11.0,
				GridSizes: []uint{2, 8, 2},
			},
			ContentProp: props.TableListContent{
				Family:    consts.Arial,
				Style:     consts.Bold,
				Size:      11.0,
				GridSizes: []uint{2, 8, 2},
			},
			Align:              consts.Center,
			HeaderContentSpace: 0.5,
		})
	})

	m.Row(5, func() {
		m.Col(12, func() {
			m.Text("OFICINA DE TECNOLOGIAS DE LA INFORMACION", props.Text{
				Size:            12,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Center,
				VerticalPadding: 0.1,
			})
		})
	})

	separatorLine(m)
	dataHojaService(m)
	separatorLine(m)
	TableStaffSuport(m)
	TableDescriptionService(m)

	contentDocument(m)
	sectionSignature(m)

	m.AddPage()
	contentDocument(m)
	//idpdf := "38be5475-6b48-4dd9-83fd-77f51dfdb97e"
	//render page in webrowser
	//bf, err := m.Output()
	//if err != nil {
	//	log.Println("could not save pdf: ", err)
	//	return
	//}
	//m.ou

	//id := uuid.New().String()
	//m.OutputFileAndClose("tmp/" + id + ".pdf")

	//w.Header().Add("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(map[string]interface{}{
	//	"status":  200,
	//	"message": "use with url api",
	//	"urls": []string{
	//		fmt.Sprintf("/temp/%s.pdf", id),
	//	},
	//})
	//os.MkdirTemp("","")
	//w.Write([]byte("ok"))
	buf, _ := m.Output()
	//image.Decode(&buf)
	buf.WriteTo(w)

	//ed := base64.StdEncoding.EncodeToString(buf.Bytes())
	//w.Header().Set("Content-Type", "text/plain")
	//w.Header().Set("Content-Transfer-Encoding", "base64")
	//w.Write([]byte(ed))
}
func info(w http.ResponseWriter, r *http.Request) {
	//os.MkdirTemp
	f, err := os.Open("tmp/38be5475-6b48-4dd9-83fd-77f51dfdb97e.pdf")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	//rd, f, err := model.NewPdfReaderFromFile("tmp/38be5475-6b48-4dd9-83fd-77f51dfdb97e.pdf", nil)
	rd, err := model.NewPdfReader(f)
	if err != nil {
		log.Println(err)
		return
	}

	v, _ := rd.GetNumPages()
	log.Println("-> ", v)
	p, _ := rd.GetPage(1)

	//render.ImageDevice
	di := render.NewImageDevice()
	di.RenderToPath(p, "tmp/hola.png")
	//i, err := di.Render(p)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println("end img render-> ", v)

	//of, err := os.CreateTemp("", "load.*.jpeg")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//defer os.Remove(of.Name())

	//log.Println("end create img-> ", v)

	//jpeg.Encode(of, i, nil)

	//w.Write([]byte(of.Name()))
	w.Write([]byte("hola"))
}

// GLOBALS DATA TABLE SUPPORT
var headerTableStaff [][]string = [][]string{
	{"1", "N°"},
	{"5", "NOMBRES Y APELLIDOS"},
	{"3", "CARGO"},
	{"3", "ASIGNACION DE SERVICIO"},
}
var bodyTableStaff []models.Person

var sizeColumBodyTableStaff []uint

func TableStaffSuport(m pdf.Maroto) {
	sizeColumBodyTableStaff = []uint{1, 5, 3, 3}
	bodyTableStaff = []models.Person{
		{FirstName: "Luis E.", LastName: "Quispe Alegre", Staff: "Analista PAD I"},
		{FirstName: "Candy Rosario", LastName: "Jara Cutipa", Staff: "Especialista en Redes"},
		{FirstName: "Jhojana", LastName: "Honorio Ferro", Staff: "Apoyo Soporte y Redes"},
		{FirstName: "Jefferson", LastName: "Morales Zavaleta", Staff: "Soporte SIGA"},
		{FirstName: "Abner", LastName: "Acuna Carrasco", Staff: "Soporte Tecnico", Assigned: true},
	}

	m.Row(2, func() {})
	makeTableStaffHeader(m)
	makeTableStaffBody(m)
	m.Row(3, func() {})

	m.SetBorder(false)
}

func makeTableStaffHeader(m pdf.Maroto) {
	m.SetBorder(true)
	m.SetBackgroundColor(getLightBlack())
	propText := props.Text{
		Size:            12,
		Family:          consts.Arial,
		Align:           consts.Center,
		Top:             2.5,
		VerticalPadding: 0.1,
		Color:           color.NewWhite(),
	}

	m.Row(10, func() {
		for index, header := range headerTableStaff {
			value, _ := strconv.Atoi(header[0])
			if index == 3 {
				m.Col(uint(value), func() {
					m.Text(header[1], props.Text{
						Size:            12,
						Family:          consts.Arial,
						Align:           consts.Center,
						Top:             0,
						VerticalPadding: 0.1,
						Color:           color.NewWhite(),
					})
				})
				continue
			}
			m.Col(uint(value), func() {
				m.Text(header[1], propText)
			})
		}
	})
	m.SetBorder(false)
}
func makeTableStaffBody(m pdf.Maroto) {
	m.SetBorder(true)
	m.SetBackgroundColor(color.NewWhite())

	var assignedStaff string = ""

	propText := props.Text{
		Size:            11,
		Family:          consts.Arial,
		Align:           consts.Center,
		VerticalPadding: 0.1,
	}

	for index, body := range bodyTableStaff {
		if body.Assigned {
			assignedStaff = "Asignado"
			m.SetBackgroundColor(getLightPurpleColor())
		}
		m.Row(5, func() {
			m.Col(sizeColumBodyTableStaff[0], func() {
				m.Text(strconv.Itoa(index+1), propText)
			})
			m.Col(sizeColumBodyTableStaff[1], func() {
				m.Text(fmt.Sprintf("%s %s", body.FirstName, body.LastName), propText)
			})
			m.Col(sizeColumBodyTableStaff[2], func() {
				m.Text(body.Staff, propText)
			})
			m.Col(sizeColumBodyTableStaff[3], func() {
				m.Text(assignedStaff, propText)
			})
		})
		m.SetBackgroundColor(color.NewWhite())
		assignedStaff = ""
	}
	m.SetBorder(false)
}

func dataHojaService(m pdf.Maroto) {
	m.Row(5, func() {
		m.Col(9, func() {
			m.Text("HOJA DE SERVICIO N°", props.Text{
				Size:            12,
				Family:          consts.Arial,
				Align:           consts.Right,
				Right:           2,
				VerticalPadding: 0.1,
			})
		})

		m.SetBorder(true)
		m.Col(3, func() {
			m.Text("001", props.Text{
				Size:            12,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Center,
				VerticalPadding: 0.1,
			})
		})
	})

	m.SetBorder(false)
	m.Row(5, func() {
		m.Col(5, func() {
			m.Text("Ciudad Universitaria", props.Text{
				Size:            12,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Center,
				VerticalPadding: 0.1,
			})
		})
		m.Col(4, func() {
			m.Text("FECHA INICIO", props.Text{
				Size:            12,
				Family:          consts.Arial,
				Align:           consts.Right,
				Right:           2,
				VerticalPadding: 0.1,
			})
		})
		m.SetBorder(true)
		m.Col(3, func() {
			m.Text("27/03/2023", props.Text{
				Size:            12,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Center,
				VerticalPadding: 0.1,
			})
		})
	})

	m.SetBorder(false)
	m.Row(5, func() {
		m.ColSpace(5)
		m.Col(4, func() {
			m.Text("FECHA FIN", props.Text{
				Size:            12,
				Family:          consts.Arial,
				Align:           consts.Right,
				Right:           8,
				VerticalPadding: 0.1,
			})
		})
		m.SetBorder(true)
		m.Col(3, func() {
			m.Text("27/03/2023", props.Text{
				Size:            12,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Center,
				VerticalPadding: 0.1,
			})
		})
		m.SetBorder(false)
	})
}

var DESCRIPTION_LIST []models.ServicesOti = []models.ServicesOti{
	{Description: "Mantenimiento de Equipos"},
	{Description: "Instalacion de Internet"},
	{Description: "Instal. o Configuracion de red"},
	{Description: "Publicacion/Pag. Web"},
	{Description: "Correos Institucionales"},
	{Description: "Instalacion de Impresora o Scaner"},
	{Description: "Help Desk"},
	{Description: "Revision o Verificacion Tecnica"},
	{Description: "Otros", Selected: true},
}

func TableDescriptionService(m pdf.Maroto) {
	m.SetBorder(true)
	m.SetBackgroundColor(getLightBlack())

	var SIZE_COLUMNS_DESCRIPTION float64 = 2.0

	columnsDescription := float64(len(DESCRIPTION_LIST)) / SIZE_COLUMNS_DESCRIPTION
	sizeRows := math.Ceil(columnsDescription)

	stylePropHeader := props.Text{
		Size:            12,
		Family:          consts.Arial,
		Align:           consts.Center,
		VerticalPadding: 0.1,
		Color:           color.NewWhite(),
	}

	m.Row(6, func() {
		m.Col(12, func() {
			m.Text("Descripcion del servicio", stylePropHeader)
		})
	})

	m.SetBackgroundColor(color.NewWhite())
	var indexDescription int = 0

	for i := 0; i < int(sizeRows); i++ {
		m.Row(5, func() {
			for j := 0; j < int(SIZE_COLUMNS_DESCRIPTION); j++ {
				if indexDescription >= len(DESCRIPTION_LIST) {
					break
				}
				if DESCRIPTION_LIST[indexDescription].Selected {
					m.SetBackgroundColor(getLightPurpleColor())
				}
				m.Col(1, func() {
					m.Text(strconv.Itoa(indexDescription+1), props.Text{
						Size:            11,
						Family:          consts.Arial,
						Align:           consts.Center,
						VerticalPadding: 0.1,
					})
				})
				m.Col(5, func() {
					m.Text(DESCRIPTION_LIST[indexDescription].Description, props.Text{
						Size:            11,
						Family:          consts.Arial,
						Align:           consts.Center,
						VerticalPadding: 0.1,
					})
				})
				m.SetBackgroundColor(color.NewWhite())
				indexDescription++
			}
		})
	}
	m.SetBackgroundColor(color.NewWhite())
}

func sectionSignature(m pdf.Maroto) {
	m.SetBorder(false)
	var signatureText []string = []string{
		"V.B del Area Tecnica",
		"Conformidad del Area o Solicitante",
	}
	//m.Row(3, func() {})
	m.RegisterFooter(func() {
		headerTableSignature := []string{
			"",
			"",
			"",
		}
		contentTableSignature := [][]string{
			{"", "", ""},
			{"", lineSignature(len(signatureText[0]) - 2), lineSignature(len(signatureText[1]) - 8)},
			{"", signatureText[0], signatureText[1]},
			{"", "luis angel pfuno canales", ""},
			{"", "72453560", ""},
		}

		m.Row(30, func() {
			m.SetBorder(true)
			m.Col(2, func() {
				m.QrCode("https://github.com/johnfercher/maroto", props.Rect{
					Left:    0,
					Top:     0,
					Center:  true,
					Percent: 90,
				})
			})
			m.SetBorder(false)
			m.ColSpace(4)
			m.TableList(headerTableSignature, contentTableSignature, props.TableList{
				HeaderProp: props.TableListContent{
					Family:    consts.Arial,
					Style:     consts.Bold,
					Size:      11.0,
					GridSizes: []uint{2, 4, 6},
				},
				ContentProp: props.TableListContent{
					Family:    consts.Arial,
					Size:      11.0,
					GridSizes: []uint{2, 4, 6},
				},
				Align:              consts.Center,
				HeaderContentSpace: 0.1,
			})
		})
	})
}
func lineSignature(sizeLine int) string {
	increment := 0
	line := ""
	for increment < (sizeLine + 6) {
		line += "_"
		increment++
	}
	return line
}

func contentDocument(m pdf.Maroto) {
	m.SetBorder(false)
	m.Row(3, func() {})
	m.Row(6, func() {
		m.Col(12, func() {
			m.Text("DETALLES U OBSERVACIONES:", props.Text{
				Size:            11,
				Style:           consts.Bold,
				Family:          consts.Arial,
				Align:           consts.Left,
				VerticalPadding: 0.1,
			})
		})
	})
}

func separatorLine(m pdf.Maroto) {
	m.Row(2, func() {})
	m.Line(2, props.Line{
		Style: consts.Solid,
		Width: 1,
	})
	m.Row(2, func() {})
}

// COLORS
func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   201,
		Green: 201,
		Blue:  201,
	}
}
func getLightBlack() color.Color {
	return color.Color{
		Red:   67,
		Green: 67,
		Blue:  71,
	}
}
