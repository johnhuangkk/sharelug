$(function () {
    $.fn.QuantityBtn = function (elm, opt) {
        let DEFAULT_OPTIONS = {};

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);
        let elem = [];
        plugin.init = function () {

            $.each($(this).find("span"), function (k, v) {
                elem[k] = $(this);
            });
            elem[0].on("click", function () {
                plugin.SpinnerDownBtn(elem);
            });
            elem[2].on("click", function () {
                plugin.SpinnerUpBtn(elem);
            });
            $("#stock").bind("change", function () {
                let count = parseInt(elem[1].html());
                let stock = parseInt($(this).val());
                if (count >= stock) {
                    count = stock;
                    elem[2].css('background-image', 'url(/static/img/icon-open-add-12-gray-70.svg)');
                } else {
                    count = 1;
                    elem[2].css('background-image', 'url(/static/img/icon-open-add-12-black-30.svg)');
                }
                elem[1].text(count).trigger('change');
            });
        }

        plugin.SpinnerDownBtn = function(element) {
            let count = element[1].html();
            if (parseInt(count) === 1) {
                count = 1;
                element[0].css('background-image', 'url(/static/img/icon-close-cut-12-gray-70.svg)');
            } else {
                count = parseInt(count) - 1;
                element[2].css('background-image', 'url(/static/img/icon-open-add-12-black-30.svg)');
            }
            element[1].text(count).trigger('change');
        };

        plugin.SpinnerUpBtn = function(element) {
            let count = element[1].html();
            let stock = $("#stock").val();
            console.log(stock);
            if (parseInt(count) >= parseInt(stock)) {
                count = stock;
                element[2].css('background-image', 'url(/static/img/icon-open-add-12-gray-70.svg)');
            } else {
                count = parseInt(count) + 1;
                element[0].css('background-image', 'url(/static/img/icon-close-cut-12-black-30.svg)');
            }
            element[1].text(count).trigger('change');
        }

        plugin.init();
    };
});