(function ($) {

    let defaults = {
        confirmButtonText: "確定",
        cancelButtonText: "取消",
        showType: "none",
        showIcon: "",
        autoClose: true,
        sensorEvent: "li",
        backdrop: false,
    };

    jQuery.fn.extend({
        halfWidth: function () {
            let width = 0;
            this.each(function () {
                width += $(this).outerWidth() / 2;
            });
            return width;
        },
        halfHeight: function () {
            let height = 0;
            this.each(function () {
                height += $(this).outerHeight() / 2;
            });
            return height;
        }
    });
    function centerWindow() {
        this._alertWindow.css({
            marginLeft: -this._alertWindow.halfWidth(),
            marginTop: -this._alertWindow.halfHeight()
        });
    }

    function createConfirmWindow(title, html) {
        //<div class='wow-alert-overlay'></div><div class='wow-alert-content'><p>" + msg + "</p><a href='#'>" + this.options.label + "</a></div>
        let elements = $('<section class="popupZone"><div class="wrap" style=""><div class="alertIconText" style="pointer-events:fill">' + this.options.showIcon +
            '<h4>' + title + '</h4>' +
            '<div class="cont">' + html + '</div>' +
            '<div class="btns">' +
            '  <span><button class="genBtnBorder h37 prime" id="cancelBtn">'+this.options.cancelButtonText+'</button></span>' +
            '  <span><button class="genBtn h37 prime" id="confirmBtn">'+this.options.confirmButtonText+'</button></span>' +
            '</div>' +
            '</div></div></section>');
        this._alertOverlay = $(elements[0]);
        this._alertWindow = $(elements[1]);
        this._actionBackdrop = this._alertOverlay.find(".wrap");
        this._actionCancelButton = this._alertOverlay.find('#cancelBtn');
        this._actionConfirmButton = this._alertOverlay.find('#confirmBtn');

        this._alertOverlay.appendTo("body");
        this._alertWindow.appendTo("body");

        return [this._alertOverlay, this._alertWindow];
    }

    function createWindow(title, html) {
        //<div class='wow-alert-overlay'></div><div class='wow-alert-content'><p>" + msg + "</p><a href='#'>" + this.options.label + "</a></div>
        let elements = $('<section class="popupZone"><div class="wrap"><div class="alertIconText">' + this.options.showIcon +
            '<h4>' + title + '</h4>' +
            '<div class="cont">' + html + '</div>' +
            '<div class="btns">' +
            '  <span><button class="genBtn h37 prime" id="confirmBtn">'+this.options.confirmButtonText+'</button></span>' +
            '</div>' +
            '</div></div></section>');
        this._alertOverlay = $(elements[0]);
        this._alertWindow = $(elements[1]);
        this._actionBackdrop = this._alertOverlay.find(".wrap");
        this._actionConfirmButton = this._alertOverlay.find('#confirmBtn');

        this._alertOverlay.appendTo("body");
        this._alertWindow.appendTo("body");

        return [this._alertOverlay, this._alertWindow];
    }


    function createPopupTimingWindow(html) {
        let elements = $('<section class="popupZone"><div class="wrap"><div class="alertIconText">'
            + this.options.showIcon + html +
            '</div></div></section>');

        this._alertOverlay = $(elements[0]);
        this._alertWindow = $(elements[1]);
        this._actionBackdrop = this._alertOverlay.find(".wrap");

        this._alertOverlay.appendTo("body");
        this._alertWindow.appendTo("body");

        setTimeout(function(){
            if (this.options.autoClose) close();
        }, 1000);

        return [this._alertOverlay, this._alertWindow];
    }

    function createAlertPopupWindow(html) {
        let elements = $('<section class="popupZone"><div class="wrap"><div><span class="close" id="closeBtn">×</span><div class="wrap2">'
            + html + '</div></div></div></section>');
        this._alertOverlay = $(elements[0]);
        this._alertWindow = $(elements[1]);
        this._actionBackdrop = this._alertOverlay.find(".wrap");
        this._actionCancelButton = this._alertOverlay.find('#closeBtn');
        this._alertOverlay.appendTo("body");
        this._alertWindow.appendTo("body");

        return [this._alertOverlay, this._alertWindow];
    }

    function createModalWindow(url) {
        let elements = $('<section class="flashList"><div class="wrap"><div><span class="close" id="closeBtn">×</span><div class="wrap2" style="overflow-y:scroll"><ul>' +
            '</ul></div></div></div></section>');
        this._alertOverlay = $(elements[0]);
        this._alertWindow = $(elements[1]);
        this._actionBackdrop = this._alertOverlay.find(".wrap");
        this._actionCancelButton = this._alertOverlay.find('#closeBtn');

        this._alertOverlay.appendTo("body");
        this._alertWindow.appendTo("body");
        let elem = $('<button class="hClose"></button>');
        this._buttonOverlay = $(elem[0]);
        this._buttonWindow = $(elem[1]);
        this._buttonOverlay.appendTo("#right");
        this._buttonWindow.appendTo("#right");
        this._actionButton = this._buttonOverlay.find(".hClose");

        let context = this;
        getData(url, this._alertOverlay.find(".wrap2"), function (){
            this._alertOverlay.find(this.options.sensorEvent).bind("click", function (e){
                e.preventDefault();
                if (context.options.response) context.options.response($(this));
                if (context.options.autoClose) close();
            });
        });

        this._alertOverlay.find(".wrap2").scroll(function (){
            console.log($(this).scrollTop() >= $(this).height());
            console.log();
        });
        $("body").css("overflow", "hidden");
        return [this._alertOverlay, this._alertWindow, this._buttonOverlay];
    }

    function configureActions() {
        let context = this;
        if (this._actionConfirmButton !== undefined) {
            this._actionConfirmButton.bind('click', function (e) {
                e.preventDefault();
                if (context.options.autoClose) close();
                if (context.options.confirm) context.options.confirm();
            });
        }
        if (this._actionCancelButton !== undefined) {
            this._actionCancelButton.bind("click", function (e) {
                e.preventDefault();
                if (context.options.autoClose) close();
                if (context.options.cancel) context.options.cancel();
            });
        }

        if (this._actionSensorButton !== undefined) {
            this._actionSensorButton.bind("click", function (e){
                e.preventDefault();
                if (context.options.response) context.options.response($(this));
                if (context.options.autoClose) close();
            });
        }

        if (this.options.backdrop) {
            this._actionBackdrop.bind("click", function (e) {
                e.preventDefault();
                if (e.target.className === "wrap") {
                    if (context.options.autoClose) close();
                }
            })
        }

        if (this._buttonOverlay !== undefined) {
            this._buttonOverlay.bind("click", function() {
                console.log("ssss");
                if (context.options.autoClose) close();
            });
        }
    }

    function close() {
        this.options.backdrop = false;
        this.options.showIcon = "";
        this._alertOverlay.remove();
        this._alertWindow.remove();
        if (this._buttonOverlay !== undefined) {
            this._buttonOverlay.remove();
        }
        $("body").css("overflow", "");
    }

    window.confirms = function (title, html, opts) {
        this.options = $.extend(defaults, opts);
        createConfirmWindow(title, html);
        centerWindow();
        configureActions();
    }


    window.alerts = function (title, html, opts) {
        this.options = $.extend(defaults, opts);
        createWindow(title, html);
        centerWindow();
        configureActions();
    }

    window.alertPopup = function (html, opts) {
        this.options = $.extend(defaults, opts);
        createAlertPopupWindow(html);
        centerWindow();
        configureActions();
    }

    window.popupTiming = function (html, opts) {
        this.options = $.extend(defaults, opts);
        createPopupTimingWindow(html);
        centerWindow();
    }

    window.modals = function (html, opts) {
        this.options = $.extend(defaults, opts);
        createModalWindow(html);
        centerWindow();
        configureActions();
    }

    function getData(url, elem, callback){
        $.http.prerequest(function () {
        }).get(url).done(function (result) {
            if (result.Status === "Success") {
                console.log(result.Data);
                let html = "";
                $.each(result.Data.ProductList, function (index, value) {
                    let spec = [];
                    $.each(value.ProductSpecList, function (i, v){
                        spec.push(v.Spec);
                    });
                    let ship = [];
                    $.each(value.ProductShipList, function (i, v){
                        ship.push(v.Text);
                    });
                    let payWay = [];
                    $.each(value.productPayWayList, function (i, v){
                        payWay.push(v.Text);
                    });

                    html += "<li data-id='" + value.ProductId + "'><p class='title'>" + value.ProductName + "</p>" +
                       "<p class='fee'><span class='price tw'>" + formatfloat(value.ProductPrice,0) + "</span></p>" +
                       "<p class='spec'>規格：" + spec.join("、") + "</p><p class='shipping'>" +
                       ship.join("、") + "</p><p class='payment'>" + payWay.join("、") + "</p></li>";
                });
                elem.find("ul").append(html);
                callback();
            } else {
                console.log(result);
            }
        });
    }

    function formatfloat(src, pos) {
        let num = parseFloat(src).toFixed(pos);
        num = num.toString().replace(/\$|\,/g,'');
        if (isNaN(num)) num = "0";
        let sign = (num !== (num = Math.abs(num)));
        num = Math.floor(num * 100 + 0.50000000001);
        let cents = num % 100;
        num = Math.floor(num/100).toString();
        if(cents<10) cents = "0" + cents;
        for (let i = 0; i < Math.floor((num.length-(1+i))/3); i++)
            num = num.substring(0, num.length-(4*i+3))+','+num.substring(num.length-(4*i+3));
        return (((sign)?'':'-') + num);
    }

}(jQuery));