$(document).ready(function() {
    $("#calc").on("click", function() {
        console.log("ran");
        var P1 = $("#P1").val();
        var P2 = $("#P2").val();
        $.ajax({
            url: "/calc",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                P1: P1,
                P2: P2,
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
    $("#clear").on("click", function() {
        $.ajax({
            url: "/clear",
            method: "GET",
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});