$(function() {
    $.fn.RealTimeReview = function (elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseProductName: $("#productName"),
            BaseProductPrice: $("#ProductPrice"),
            BaseShipText: $("#ShipText"),
            BaseShipFee: $("#ShipFee"),
            BasePayWayText: $("#PayWayText"),
            BaseTotal: $("#Total"),
            BaseQrCodeImage: $(".qrcode"),
            BaseShortUrl: $("#ShortUrl"),
            BaseUpLoadImg: $(".uploadImgS"),
            BaseSeller: $(".seller"),
            BaseData:"",
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        plugin.init = function (){
            plugin.view();

            $(".copy").on("click", function () {
                popupTiming("<p>連結已複製到剪貼簿！</p>", {
                    showIcon: '<img class="m10" src="/static/img/icon-checkL.svg">',
                });
                let temp = $('<input>'); // 建立input物件
                $('body').append(temp); // 將input物件增加到body
                let url = $("#short").val(); // 取得要複製的連結
                temp.val(url).select(); // 將連結加到input物件value
                document.execCommand('copy'); // 複製
                temp.remove(); // 移除input物件
            });

            $(".showQrcode").bind("click", function () {
                plugin.QrCodeAlert();
            });

            $(".shareAlert").on("click", function () {
                plugin.shareAlert();
            });
        };

        plugin.view = function () {
            // console.log(options.BaseData);
            let shipfee = 0;

            options.BaseProductName.text(options.BaseData.ProductName);
            options.BaseProductPrice.text(options.BaseData.Price);

            $.each(options.BaseData.ShippingList, function (index, value) {
                options.BaseShipText.text(value.Text);
                if (value.Type === "none") {
                    options.BaseShipFee.removeClass("tw");
                    options.BaseShipFee.text(value.Remark);
                } else {
                    options.BaseShipFee.text($(this).formatfloat(value.Price));
                }
                shipfee = value.Price;
            });

            options.BaseTotal.text($(this).formatfloat(options.BaseData.Price + shipfee));

            $.each(options.BaseData.ProductPayWayList, function (index, value) {
                options.BasePayWayText.text(value.Text);
            });

            options.BaseQrCodeImage.find("img").attr("src", options.BaseData.QrCode);
            options.BaseQrCodeImage.find("a").attr("href", options.BaseData.QrCode);

            $("#short").val(options.BaseData.ShortUrl);
            options.BaseShortUrl.text(options.BaseData.ShortUrl);
            // console.log(options.BaseData.ExpireDate);

            $('#clock1').countdown(options.BaseData.ExpireDate, function(event) {
                let totalHours = event.offset.totalDays * 24 + event.offset.hours;
                $(this).html(event.strftime(totalHours + ':%M:%S'));
            });
            $('#clock2').countdown(options.BaseData.ExpireDate, function(event) {
                let totalHours = event.offset.totalDays * 24 + event.offset.hours;
                $(this).html(event.strftime(totalHours + ':%M:%S'));
            });

            if (options.BaseData.ProductImageList !== null) {
                $.each(options.BaseData.ProductImageList, function (index, value) {
                    options.BaseUpLoadImg.find("label").css("background-image", "url(" + value + ")");
                });
            } else {
                options.BaseUpLoadImg.hide();
            }
            $("#footer").addClass("hideM");
            $(".wrapper").addClass("wh50");

        };

        plugin.QrCodeAlert = function () {
            alertPopup("<img class='qrcodeP' src='"+options.BaseData.QrCode+"' />", {
                backdrop:true
            });
        }

        plugin.shareAlert = function () {
            let $dialog = $("#specModal");
            $dialog.find(".close").on("click",function (){
                $dialog.hide();
            });
            $dialog.show();
        }

        plugin.init();
    }
});
