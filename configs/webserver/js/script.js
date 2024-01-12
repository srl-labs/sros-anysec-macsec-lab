$(function () {
    $('.main-container .left a').on('click', function (e) {
        e.preventDefault();
        id = e.currentTarget.id;
        path = id.replace(/-/g, '/');
        message = "";
        if (id.startsWith('-icmp')) {
            icmp_size = $('#icmp-size')[0].value;
            if (icmp_size > 8000) {
                message += "The ICMP size range must be between 0 and 8000 bytes\n";
                $('#icmp-size')[0].value = "8000";
            }

            if (icmp_size < 0) {
                message += "The ICMP size range must be between 0 and 8000 bytes\n";
                $('#icmp-size')[0].value = "0";
            }
            path += '/' + icmp_size
            icmp_interval = $('#icmp-interval')[0].value;
            if (icmp_interval > 1) {
                message += "The ICMP interval range must be between 0.01 and 1 second\n";
                $('#icmp-interval')[0].value = "1";
            }

            if (icmp_interval < 0.01) {
                message += "The ICMP interval range must be between 0.01 and 1 second\n";
                $('#icmp-interval')[0].value = "0.01";
            }
            path += '/' + icmp_interval
        }
        console.log(path);
        if (message != "") {
            alert(message);
        }
        else {
            $.get(path,
                function (data) {
                    //do nothing
                    console.log(data);
                }
            );
        }
        return false;
    });
});

function update_link_buttons() {
    path = '/get_links/data';
    $.get(path,
        function (data) {
            //do nothing
            //console.log(data);
            Object.keys(data).forEach(key => {
                if ($('.link-' + key).hasClass('disabled')) {
                    $('.link-' + key).removeClass('disabled');
                }
                if ($('.link-' + key).hasClass('enabled')) {
                    $('.link-' + key).removeClass('enabled');
                }
                $('.link-' + key).addClass(data[key]);
            });

        }
    );
}

function update_anysec_buttons() {
    path = '/get_anysecs/data';
    $.get(path,
        function (data) {
            //do nothing
            //console.log(data);
            Object.keys(data).forEach(key => {
                if ($('.anysec-' + key).hasClass('disabled')) {
                    $('.anysec-' + key).removeClass('disabled');
                }
                if ($('.anysec-' + key).hasClass('enabled')) {
                    $('.anysec-' + key).removeClass('enabled');
                }
                $('.anysec-' + key).addClass(data[key]);
            });

        }
    );
}
function update_icmp_buttons() {
    path = '/get_icmp_status';
    $.get(path,
        function (data) {
            //do nothing
            //console.log(data);
            Object.keys(data).forEach(key => {
                if ($('.icmp-' + key).hasClass('disabled')) {
                    $('.icmp-' + key).removeClass('disabled');
                }
                if ($('.icmp-' + key).hasClass('enabled')) {
                    $('.icmp-' + key).removeClass('enabled');
                }
                $('.icmp-' + key).addClass(data[key]);
            });

        }
    );
}
$(document).ready(function () {
    setInterval(update_link_buttons, 1000);
    setInterval(update_anysec_buttons, 1000);
    setInterval(update_icmp_buttons, 1000);
});