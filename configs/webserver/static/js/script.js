$(function () {
    $('.card button').on('click', function (e) {
        e.preventDefault();
        el = $(e.currentTarget)
        form = el.closest('form');
        id = e.currentTarget.id;
        path = id.replace(/-/g, '/');
        makeRequest = true;
        message = "";
        if (id.startsWith('-icmp')) {
            icmp_size = $('#icmp-size')[0].value;
            icmp_interval = $('#icmp-interval')[0].value;
            path += '/' + icmp_size + '/' + icmp_interval
        }
        console.log(path);
        if (form.eq(0).hasClass('needs-validation'))
            if (!form[0].checkValidity()) {
                e.stopPropagation();
                makeRequest = false;
            }
        if (makeRequest) {
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
            Object.keys(data).forEach(key => {
                if (data[key] == 'enabled') {
                    if ($('.link-' + key).hasClass('btn-primary')) {
                        $('.link-' + key).removeClass('btn-primary');
                    }
                    if ($('.link-' + key).hasClass('btn-danger')) {
                        $('.link-' + key).removeClass('btn-danger');
                    }
                    $('.link-' + key).addClass('btn-success');
                }
                else if (data[key] == 'disabled') {
                    if ($('.link-' + key).hasClass('btn-primary')) {
                        $('.link-' + key).removeClass('btn-primary');
                    }
                    if ($('.link-' + key).hasClass('btn-success')) {
                        $('.link-' + key).removeClass('btn-success');
                    }
                    $('.link-' + key).addClass('btn-danger');
                }


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
            Object.keys(data).forEach(key => {
                if (data[key] == 'enabled') {
                    if ($('.anysec-' + key).hasClass('btn-primary')) {
                        $('.anysec-' + key).removeClass('btn-primary');
                    }
                    if ($('.anysec-' + key).hasClass('btn-danger')) {
                        $('.anysec-' + key).removeClass('btn-danger');
                    }
                    $('.anysec-' + key).addClass('btn-success');
                }
                else if (data[key] == 'disabled') {
                    if ($('.anysec-' + key).hasClass('btn-primary')) {
                        $('.anysec-' + key).removeClass('btn-primary');
                    }
                    if ($('.anysec-' + key).hasClass('btn-success')) {
                        $('.anysec-' + key).removeClass('btn-success');
                    }
                    $('.anysec-' + key).addClass('btn-danger');
                }
            });

        }
    );
}
function update_icmp_buttons() {
    path = '/get_icmp_status';
    $.get(path,
        function (data) {
            Object.keys(data).forEach(key => {
                if (data[key] == 'enabled') {
                    if ($('.icmp-' + key).hasClass('btn-primary')) {
                        $('.icmp-' + key).removeClass('btn-primary');
                    }
                    if ($('.icmp-' + key).hasClass('btn-danger')) {
                        $('.icmp-' + key).removeClass('btn-danger');
                    }
                    $('.icmp-' + key).addClass('btn-success');
                }
                else if (data[key] == 'disabled') {
                    if ($('.icmp-' + key).hasClass('btn-primary')) {
                        $('.icmp-' + key).removeClass('btn-primary');
                    }
                    if ($('.icmp-' + key).hasClass('btn-success')) {
                        $('.icmp-' + key).removeClass('btn-success');
                    }
                    $('.icmp-' + key).addClass('btn-danger');
                }
            });

        }
    );
}
$(document).ready(function () {
    setInterval(update_link_buttons, 1000);
    setInterval(update_anysec_buttons, 1000);
    setInterval(update_icmp_buttons, 1000);
});
