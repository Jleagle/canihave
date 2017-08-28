
var loading = false;

var t = setInterval(checkPagePosition, 1000);
// clearInterval(t);

$(window).scroll(checkPagePosition);

function checkPagePosition() {
    if ($(window).scrollTop() + $(window).height() > $(document).height() - 1000) {
        if (!loading) {
            loadPage(1);
        }
    }
}

function loadPage(page) {
    $.ajax({
        dataType: "html",
        beforeSend: function () {
            loading = true;
        },
        xhr: function () {
            var xhr = $.ajaxSettings.xhr();
            
            xhr.addEventListener("progress", function (evt) {
                if (evt.lengthComputable) {
                    var percentComplete = evt.loaded / evt.total * 100;
                    console.log(percentComplete);

                    if (percentComplete == 100) {
                        $("#loading").hide();
                    }
                    else {
                        $("#loading").show();
                        $("#loading .progress-bar").width(percentComplete + "%");
                    }
                }
            }, false);

            return xhr;
        },
        type: 'GET',
        url: "/ajax",
        success: function (data, textStatus, jqXHR) {
            $("#results").append(data);
            loading = false;
        }
    });
}
