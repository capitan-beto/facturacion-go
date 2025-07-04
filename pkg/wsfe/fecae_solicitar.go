package wsfe

import (
	"cmd/api/main.go/pkg/utils"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Envelope struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	XMLNSs    string   `xml:"xmlns:soapenv,attr"`
	XMLNSwsaa string   `xml:"xmlns:ar,attr"`
	Header    string   `xml:"soapenv:Header"`
	Body      *Body    `xml:"soapenv:Body"`
}

type Body struct {
	FECAESolicitar *FECAESolicitarBody `xml:"ar:FECAESolicitar"`
}

type FECAESolicitarBody struct {
	Auth     *AuthBody `xml:"ar:Auth"`
	FeCAEReq *FeCAEReq `xml:"ar:FeCAEReq"`
}

//children of FECAESolicitar

// Información de auth. Obligatorio.
type AuthBody struct {
	// Token devuelto por wsaa. Obligatorio.
	Token string `xml:"ar:Token"`

	// Sign devuelto por wsaa. Obligatorio.
	Sign string `xml:"ar:Sign"`

	// Cuit del representado o emisor. Obligatorio.
	Cuit int64 `xml:"ar:Cuit"`
}

// Informacion del comprobante o lote de comprobantes de ingreso. Obligatorio.
type FeCAEReq struct {
	// Información de la cabecera del comprobante. Obligatorio.
	FeCabReq *FeCabReqBody `xml:"ar:FeCabReq"`

	//Informacion del detalle del comprobante. Obligatorio.
	FeDetReq *FeDetReqBody `xml:"ar:FeDetReq"`
}

// Campos de la cabecera del comprobante.
type FeCabReqBody struct {
	// Cantidad del registro del detalle del comprobante. Obligatorio.
	CantReg int `xml:"ar:CantReg"`
	// Punto de venta del comprobante que se está informando. Obligatorio.
	PtoVta int `xml:"ar:PtoVta"`
	// Tipo de comprobante que se esta informando. Obligatorio.
	CbteTipo int `xml:"ar:CbteTipo"`
}

// Campos del detalle del comprobante.
type FeDetReqBody struct {
	FECAEDetRequest *FECAEDetRequestBody `xml:"ar:FECAEDetRequest"`
}

type FECAEDetRequestBody struct {
	// Concepto del comprobante. Valores permitidos,
	//1-Productos, 2-Servicios, 3-Productos y servicios. Obligatorio.
	Concepto int `xml:"ar:Concepto"`

	// Código de documento identificatorio del comprador. Obligatorio.
	DocTipo int8 `xml:"ar:DocTipo"`

	//Nro. de identificación del comprador. Obligatorio.
	DocNro int64 `xml:"ar:DocNro"`

	//Nro de comprobante desde 1-99999999. Obligatorio.
	CbteDesde int32 `xml:"ar:CbteDesde"`

	//Nro de comprobante registrado hasta 1-99999999. Obligatorio.
	CbteHasta int32 `xml:"ar:CbteHasta"`

	//Fecha del comprobante (yyyymmdd). Para concepto igual a 1la fecha de
	//emisión del comprobante puede ser hasta 5 días anteriores o
	//posteriores respecto de la fecha de generación. No obligatorio.
	CbteFch string `xml:"ar:CbteFch,omitempty"`

	//Importe total del comprobante. Obligatorio.
	ImpTotal float64 `xml:"ar:ImpTotal"`

	//Importe neto no gravado. Obligatorio.
	ImpTotConc float64 `xml:"ar:ImpTotConc"`

	//Importe neto gravado. Debe ser menor o
	//igual a Importe total y no puede ser menor a cero. Obligatorio.
	ImpNeto float64 `xml:"ar:ImpNeto"`

	//Importe excento. Debe ser menor o igual a
	//Importe total y no puede ser menor a cero. Obligatorio.
	ImpOpEx float64 `xml:"ar:ImpOpEx"`

	//Suma de los conceptos del array de IVA. Obligatorio.
	ImpIVA float64 `xml:"ar:ImpIVA"`

	//Suma de los importes del array de tributos. Obligatorio.
	ImpTrib int64 `xml:"ar:ImpTrib"`

	//Fecha de inicio del abono para el servicio a facturar. No obligatorio.
	FchServDesde string `xml:"ar:FchServDesde,omitempty"`

	//Fecha del fin del abono para el servicio a facturar. No obligatorio.
	FchServHasta string `xml:"ar:FchServHasta,omitempty"`

	//Fecha del vencimiento del pago servicio a facturar. No obligatorio.
	FchVtoPago string `xml:"ar:FchVtoPago,omitempty"`

	//Codigo de moneda de comprobante. Obligatorio.
	MonId string `xml:"ar:MonId"`

	//Cotizacion de la moneda informada. Peso debe ser 1.  Obligatorio.
	MonCotiz float32 `xml:"ar:MonCotiz"`

	//Marca que identifica si el comprobante se cancela en misma moneda del comp.
	//No obligatorio.
	CanMisMonExt string `xml:"ar:CanMisMonExt,omitempty"`

	//Condición frente al IVA del receptor. Consultar método
	//"GET Tipo de Responsables". No obligatorio.
	CondicionIVAReceptorId int8 `xml:"ar:CondicionIVAReceptorId,omitempty"`

	//Array para informar los comprobantes asociados <CbteAsoc>. No obligatorio.
	CbtesAsoc *CbtesAsoc `xml:"ar:CbtesAsoc,omitempty"`

	//Array para informar los tributos asociados a un comprobante <Tributo>.
	//No obligatorio.
	Tributos *Tributos `xml:"ar:Tributos,omitempty"`

	//Array para informar las alícuotas y sus importes asociados a un comprobante
	//<AlicIva>. No obligatorio.
	Iva *Iva `xml:"ar:Iva,omitempty"`

	//Array de campos auxiliares. Reservado usos futuros <Opcional>.
	//Adicionales por R.G. No obligatorio.
	Opcionales *Opcionales `xml:"ar:Opcionales,omitempty"`

	//Array para informar los múltiples compradores. No obligatorio.
	Compradores *Compradores `xml:"ar:Compradores,omitempty"`

	//Estructura compuesta por la fecha desde y la fecha hasta del periodo
	//que se quiere identificar. No obligatorio.
	PeriodoAsoc *PeriodoAsoc `xml:"ar:PeriodoAsoc,omitempty"`

	//Array para informar las actividades asociadas a un comprobante.
	//No obligatorio.
	Actividades *Actividades `xml:"ar:Actividades,omitempty"`
}

// Detalle de los comprobantes relacionados con el comprobante que se
// solicita autorizar (array).
type CbtesAsoc struct {
	CbteAsoc *[]CbteAsocBody `xml:"ar:CbteAsoc"`
}

// Detalle de tributos relacionados con el comprobante que se solicita
// autorizar (array).
type Tributos struct {
	Tributo *[]TributoBody `xml:"ar:Tributo"`
}

// Detalle de alícuotas relacionadas con el comprobante que se solicita autorizar (array).
type Iva struct {
	AlicIva *[]AlicIvaBody `xml:"ar:AlicIva"`
}

// Los datos opcionales sólo deberán ser incluidos si el emisor pertenece al conjunto de emisores
// habilitados a informar opcionales. Consultar docs wsfev1. (array).
type Opcionales struct {
	Opcional *[]OpcionalBody `xml:"ar:Opcional"`
}

// Detalle compradores vinculados al comprobante que se solicita autorizar (array).
type Compradores struct {
	Comprador *[]CompradorBody `xml:"ar:Comprador"`
}

// Estructura que permite soportar un rango de fechas.
type PeriodoAsoc struct {
	//Fecha correspondiente al inicio del periodo de los comprobantes que se
	//quiere identificar. Obligatorio.
	FchDesde string `xml:"ar:FchDesde"`

	//Fecha correspondiente al fin del periodo de los comprobantes que se quiere
	//identificar. Obligatorio
	FchHasta string `xml:"ar:FchHasta"`
}

// Detalle de la actividad relacionada con las actividades (array) que se
// indican en el comprobante a autorizar.
type Actividades struct {
	Actividad *[]ActividadBody `xml:"ar:Actividad"`
}

type CbteAsocBody struct {
	//Codigo de tipo de comprobante. Consultar metodo "FEParamGetTiposCbte".
	//Obligatorio.
	Tipo int16 `xml:"ar:Id"`

	//Punto de venta del comprobante asociado. Obligatorio.
	PtoVenta int32 `xml:"ar:Desc"`

	// Numero del comprobante asociado. Obligatorio
	Nro int32 `xml:"ar:Alic"`

	// Cuit emisor del comprobante asociado. No obligatorio.
	Cuit string `xml:"ar:BaseImp,omitempty"`

	//Fecha del comprobante asociado. No obligatorio.
	CbteFch string `xml:"ar:Importe,omitempty"`
}

type TributoBody struct {
	// Código tributo según método FEParamGetTiposTributos. Obligatorio
	Id int8 `xml:"ar:Id"`

	//Descripción del tributo. No obligatorio.
	Desc string `xml:"ar:Desc,omitempty"`

	//Base imponible para la determinación del tributo. Obligatorio.
	BaseImp float64 `xml:"ar:BaseImp"`

	//Alicuota. Obligatorio.
	Alic float32 `xml:"ar:Alic"`

	//Importe del tributo. Obligatorio.
	Importe float32 `xml:"ar:Importe"`
}

type AlicIvaBody struct {
	//Código de tipo de iva. Consultar método FEParamGetTiposIva. Obligatorio.
	Id int16 `xml:"ar:Id"`

	//Base imponible para la determinación de la alícuota. Obligatorio.
	BaseImp float64 `xml:"ar:BaseImp"`

	//Importe. Obligatorio.
	Importe float64 `xml:"ar:Importe"`
}

type OpcionalBody struct {
	//Código de Opcional, consultar método FEParamGetTiposOpcional. Obligatorio.
	Id string `xml:"ar:Id"`

	//Valor. Obligatorio.
	Valor string `xml:"ar:Valor"`
}

type CompradorBody struct {
	//Tipo de documento del comprador. Obligatorio.
	DocTipo int8 `xml:"ar:DocTipo"`

	//Número de documento del comprador. Obligatorio.
	DocNro int64 `xml:"ar:DocNro"`

	//Porcentaje de titularidad que tiene el comprador. Obligatorio.
	Porcentaje float32 `xml:"ar:Porcentaje"`
}

type ActividadBody struct {
	//Código actividad según método FEParamGetActividades. Obligatorio.
	Id int64 `xml:"ar:Id"`
}

// func GetImporteTotal()

func FECAESolicitar() {
	soapenv := "http://schemas.xmlsoap.org/soap/envelope/"
	ar := "http://ar.gov.afip.dif.FEV1/"

	// Orden de creacion de xml
	// 1 - Información de cabecera.
	newFeCabReqBody := FeCabReqBody{CantReg: 1, PtoVta: 1, CbteTipo: 1}
	// 2 - Cuerpo de la solicitud de comprobante.
	newFeCAEDetRequestBody := FECAEDetRequestBody{
		Concepto:     1,
		DocTipo:      80,
		DocNro:       20356684320,
		CbteDesde:    1,
		CbteHasta:    1,
		CbteFch:      "20250205",
		ImpTotal:     12100,
		ImpTotConc:   0,
		ImpNeto:      10000,
		ImpOpEx:      0,
		ImpTrib:      0,
		ImpIVA:       2100,
		MonId:        "PES",
		MonCotiz:     1,
		CanMisMonExt: "N",
		Iva:          &Iva{AlicIva: &[]AlicIvaBody{{Id: 5, BaseImp: 10000, Importe: 2100}}},
	}

	// 3 - Anidar los detalles del comprobante al campo de cuerpo de la solicitud.
	newFeDetReqBody := FeDetReqBody{&newFeCAEDetRequestBody}

	// 4 - Anidar cabecera y detalles del comprobante información del/los comprobantes de ingreso.
	newFeCAEReq := FeCAEReq{FeCabReq: &newFeCabReqBody, FeDetReq: &newFeDetReqBody}

	// 5- Anidar los datos de autorización.
	reqAuth := AuthBody{
		Token: "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9InllcyI/Pgo8c3NvIHZlcnNpb249IjIuMCI+CiAgICA8aWQgc3JjPSJDTj13c2FhaG9tbywgTz1BRklQLCBDPUFSLCBTRVJJQUxOVU1CRVI9Q1VJVCAzMzY5MzQ1MDIzOSIgZHN0PSJDTj13c2ZlLCBPPUFGSVAsIEM9QVIiIHVuaXF1ZV9pZD0iNzQ5ODUxNDE2IiBnZW5fdGltZT0iMTczODYzODk4MCIgZXhwX3RpbWU9IjE3Mzg2ODIyNDAiLz4KICAgIDxvcGVyYXRpb24gdHlwZT0ibG9naW4iIHZhbHVlPSJncmFudGVkIj4KICAgICAgICA8bG9naW4gZW50aXR5PSIzMzY5MzQ1MDIzOSIgc2VydmljZT0id3NmZSIgdWlkPSJTRVJJQUxOVU1CRVI9Q1VJVCAyMDQwNjE2NTI4MSwgQ049ZmFjdHVyYWRvcnRlc3QxIiBhdXRobWV0aG9kPSJjbXMiIHJlZ21ldGhvZD0iMjIiPgogICAgICAgICAgICA8cmVsYXRpb25zPgogICAgICAgICAgICAgICAgPHJlbGF0aW9uIGtleT0iMjA0MDYxNjUyODEiIHJlbHR5cGU9IjQiLz4KICAgICAgICAgICAgPC9yZWxhdGlvbnM+CiAgICAgICAgPC9sb2dpbj4KICAgIDwvb3BlcmF0aW9uPgo8L3Nzbz4K",
		Sign:  "Jh8igUl/2L3FlJTFxVtWXMDZ77QP4FHbnWaVHdJIJDjF/IoP/9j1Jgvzr0NGSNuDPaK2Fkp4sJ3KVTjc11/jFgRR3izbNwasRf2X6LQchYr2LcVzCh6t82NKBWuL5Wt0mxC2JiiPik7f8rM3yKoo+kSK3hDo/JCtmMoET/Ol7Kg=",
		Cuit:  20406165281,
	}

	// Anidar datos de auth y de comprobante de ingreso al cuerpo de la solicitud.
	reqBody := FECAESolicitarBody{Auth: &reqAuth, FeCAEReq: &newFeCAEReq}
	// Añadir solicitud al cuerpo del sobre.
	newBody := Body{FECAESolicitar: &reqBody}
	// Añadir cuerpo del sobre al resto de la petición.
	newReq := Envelope{XMLNSs: soapenv, XMLNSwsaa: ar, Body: &newBody}

	out, _ := xml.MarshalIndent(newReq, " ", "  ")
	res := utils.EscapeXML(string(out))

	fmt.Println(res)

	req, err := http.NewRequest("POST", "https://wswhomo.afip.gov.ar/wsfev1/service.asmx", strings.NewReader(res))
	if err != nil {
		logrus.Error(err)
		return
	}

	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECAESolicitar")

	transport := &http.Transport{
		Proxy: nil,
	}

	client := &http.Client{
		Transport: transport,
	}

	newRes, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	body, err := io.ReadAll(newRes.Body)
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Println(string(body))

}
