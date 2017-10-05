$(window).resize(setSquareHeight);
$(document).ready(setSquareHeight);

function setSquareHeight() {
    var $squares = $('div.card a.square');
    var height = $squares.eq(0).width();

    if ($squares.length > 0) {
        $squares.height(height)
    }
}
