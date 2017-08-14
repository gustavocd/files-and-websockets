const file = document.querySelector('#file')
const ws = new WebSocket(`ws://${window.location.host}/ws/upload`)
const canvas = document.querySelector('canvas');

file.addEventListener('change', handleFile);

function handleFile(e) {
    let file = e.target.files[0]
    if (!file) return;
    compress(file)
}

function compress(file) {
    new ImageCompressor(file, {
        quality: 0.6,
        success(result) {
            upload(result)
        },
        error(e) {
            displayError(e.message);
        }
    });
}

function upload(file) {
    let reader = new FileReader();
    let res = reader.onload = function() { 
        let data = { file: reader.result, filter: makeid() };
        console.log(data);
        ws.send(JSON.stringify(data));
    };
    reader.readAsDataURL(file);
    //let blob = new Blob([data]);
}


function displayError(error) {
    alert(error);
}

ws.addEventListener('message', e => {
    let imageURL = URL.createObjectURL(e.data);
    document.querySelector('#image').classList.add('isActive');
    document.querySelector('#image').src = imageURL;
});

ws.addEventListener('close', e => {
    //ws.close();
});

function makeid() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    for (var i = 0; i < 5; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    return text;
}