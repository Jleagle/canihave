
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
    var $results = $('#results');

    if ($results.attr('data-last') != '0') {
        return
    }

    var nextPage = +$results.attr('data-page') + 1;
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
        data: {
            page: nextPage,
            search: window.atob($results.attr('data-search')),
        },
        url: "/ajax",
        success: function (data, textStatus, jqXHR) {
            data = $.trim(data);
            $data = $(data).attr('data-page', nextPage)

            if (data.length > 0) {
                $results.append($data);
                $results.attr('data-page', nextPage)
            } else {
                $results.attr('data-last', 1)
            }
            loading = false;
        }
    });
}
