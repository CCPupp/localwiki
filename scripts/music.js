var $fileInput = $("#file-input");
var $btn = $("button");
var isPlaying = false;
var player;

function readFile(e) {
    if (window.FileReader) {
        var file = e.target.files[0];
        var reader = new FileReader();
        if (file && file.type.match("audio/flac")) {
            reader.readAsArrayBuffer(file);
        } else {
            console.log("Please add a flac file.");
        }
        reader.onloadend = function(e) {
            player = AV.Player.fromBuffer(reader.result);
            $btn.show();
            $btn.on("click", function() {
                if (isPlaying) {
                    player.pause();
                    isPlaying = false;
                    this.textContent = "play";
                } else {
                    player.play();
                    isPlaying = true;
                    this.textContent = "pause";
                }
            });
        };
    }
}

$fileInput.on("change", readFile);