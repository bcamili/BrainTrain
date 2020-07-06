var simplemde1 = new SimpleMDE({
  element: document.getElementById("mde1"),
  placeholder: "Gib eine Frage ein",
});
var simplemde2 = new SimpleMDE({
  element: document.getElementById("mde2"),
  placeholder: "Gib die richtige Antwort an",
});

var simplemde3 = new SimpleMDE({
  element: document.getElementById("mde3"),
});

var simplemde4 = new SimpleMDE({
  element: document.getElementById("mde4"),
});

var header = document.getElementById("navigationButtons");
var btns = header.getElementsByTagName("a");
for (var i = 0; i < btns.length; i++) {
  if (document.location.href.indexOf(btns[i].href) >= 0) {
    btns[i].className += " active";
  }
};

function on(id) {

  if(id!='overlay'){
    console.log(id)
  id = id.substring(1, 33)
}

  document.getElementById(id).style.display = "block";
};

function off(id) {
  if(id!='overlay'){
    console.log(id)
  id = id.substring(1, 33)
}
console.log(id)
  document.getElementById(id).style.display = "none";
};

function onCard(id) {
id = id.substring(1, 81)
console.log(id)

  document.getElementById(id).style.display = "block";
};

function offCard(id) {
  id = id.substring(1, 81)
console.log(id)
  document.getElementById(id).style.display = "none";
};

function fileUploader() {
  document.getElementById("overlayEdit").style.display = "block";
};

function fileUploaderOff() {
  document.getElementById("overlayEdit").style.display = "none";
};
