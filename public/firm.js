//imported in main.js
//const fxHoja = document.getElementById("fxhoja")
const fxcontent = document.getElementById("fxcontent")
const fxfirm = document.getElementById("fxfirm")
let move = false

$(document).ready(function(){
  fxhoja.addEventListener('mousemove',e => {
    if(move){mover_a(e)}
  })
  fxhoja.addEventListener('mousedown',()=> {
    move = true
  })

  //move when press container
  document.addEventListener("mouseup",()=>{
    move = false
  })
})



function mover_a(event) {
  const relative_pos = get_relative(event);
  var new_left = relative_pos.x - (fxfirm.offsetWidth / 2);
  var new_top = relative_pos.y - (fxfirm.offsetHeight / 2);

  if(new_left < 0 ){
      new_left = 0;
  } else if (new_left > fxhoja.offsetWidth - fxfirm.offsetWidth){
      new_left = fxhoja.offsetWidth - fxfirm.offsetWidth;
  }

  if(new_top < 0 ){
      new_top = 0;
  } else if(new_top > fxhoja.offsetHeight - fxfirm.offsetHeight){
      new_top = fxhoja.offsetHeight - fxfirm.offsetHeight;
  }

  fxfirm.style.left = new_left + 'px';
  fxfirm.style.top = new_top + 'px';
}
function get_relative(event) {
    const pos = event.currentTarget.getBoundingClientRect();
    return {
        x: event.clientX - pos.left,
        y: event.clientY - pos.top
    };
}


//send button init firm event initInvoker
function send_firm() {
  initInvoker('W')
}

//getArgsForEvent send request to api to get info base64
const getArgsForEvent =async(payload)=>{
  const rq = await fetch("http://localhost:4000/args",{
    method:"POST",
    headers:{
      "Content-Type":"application/json",
    },
    body:JSON.stringify(payload),
  })

  const body = await rq.json()
  return body
  //console.log(body)
}

//events implement required
window.addEventListener('getArguments',async(e)=>{
  type = e.detail;
  let position = $("#fxfirm").position()
  let payload = {
    file_id:"38be5475-6b48-4dd9-83fd-77f51dfdb97e",
    page_number:"0",
    exacto:1,
    poy:position.top.toFixed(0),
    pox:position.left.toFixed(0),
    stamp_appearance_id:"0",
    reason:"Soy el autor del documento",
  }
  console.log(payload.pox," == ",payload.poy)
  let args = await getArgsForEvent(payload)
  console.log(args)
  dispatchEventClient('sendArguments', args.args);
})
window.addEventListener('invokerOk',()=>{
  console.log("todo ok ya firmo")
})
window.addEventListener('invokerCancel',(e)=>{
  alerta("El proceso de firma digital fue cancelado.", false);
})
