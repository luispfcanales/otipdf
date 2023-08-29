//url
//let url = "http://localhost:4000/pdf"
const fxHoja = document.getElementById("fxhoja")

let pdfjsLib = window['pdfjs-dist/build/pdf'];

pdfjsLib.GlobalWorkerOptions.workerSrc = '//mozilla.github.io/pdf.js/build/pdf.worker.js';

//let loadingtask = pdfjsLib.getDocument("https://unamadpdf.onrender.com/pdf")
let loadingtask = pdfjsLib.getDocument("http://18.118.181.184/pdf")

const factor = 2.5
const datos_firma = {
  archivo_id:0,
  num_pagina:0,
  motivo:'soy el autor',
  exacto:1,
  pos_pagina:"0-0",
  apariencia:0
}

function printpdf() {
  //const urlblob = URL.createObjectURL(vblob)
  //const url = `${urlblob}.pdf`

  loadingtask.promise.then(function(pdf){
    pdf.getPage(1).then(function(page){
      let viewport = page.getViewport({scale:0.88})
      let canvas = document.createElement('canvas')
      let ctx = canvas.getContext('2d')
      canvas.width = 524.92;
      canvas.height = 742.41;

      page.render({canvasContext:ctx,viewport:viewport}).promise.then(function(){
        const imgcanvas = canvas.toDataURL('image/png')
        fxHoja.style.position = "relative";
        fxHoja.style.width = canvas.width+"px";
        fxHoja.style.height = canvas.height+"px";
        fxHoja.style.backgroundImage = `url("${imgcanvas}")`
      })
    })
  })

}

async function LoadPdf() {
  //const r = await fetch(url)//.then(response => response.atob())
  //const value = await r.blob()
  printpdf()
}
LoadPdf()


