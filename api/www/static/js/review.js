$(function() {
    $.fn.ProductReview = function (elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseImages: $("#product_image"),
            BaseImageBtn: $("#product_image_btn"),
            BaseName: $("#product_name"),
            BasePrice: $("#price"),
            BaseQRCodeLink: $(".qrcode a"),
            BaseQRCodeImage: $(".qrcode img"),
            BaseShortUrl: $("#ShortUrl"),
            BaseSpecList: $(".spec .select"),
            BaseShipList: $("#shipping"),
            BaseSpec: $("#spec"),
            BaseCountnum: $("#countnum"),
            BaseTotal: $("#total"),
            BaseDownBtn: $(".down"),
            BaseUpBtn: $(".up"),
            BaseData: "",
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        plugin.init = function() {
            // console.log(options.BaseData);
            plugin.view(options.BaseData);

            options.BaseSpecList.find("select").on("change", function () {
                let price = options.BaseData.Price;
                let Shipping = options.BaseShipList.find("option:selected").data("fee");
                let quantity = options.BaseCountnum.html();
                // console.log(price , quantity , Shipping)
                options.BaseTotal.text($(this).formatfloat(price * quantity + Shipping, 0));
            });

            options.BaseDownBtn.click(function () {
                plugin.SpinnerDownBtn();
            });
            options.BaseUpBtn.click(function () {
                plugin.SpinnerUpBtn();
            });

            options.BaseShipList.on("change", function () {
                plugin.CalculateTotal();
            });

            options.BaseCountnum.on("change", function () {
                plugin.CalculateTotal();
            });

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

        plugin.view = function (BaseData) {
            //圖片
            $.each(BaseData.ProductImageList, function (k, v) {
                options.BaseImages.append("<li style=\"background-image: url(" + v + ")\"></li>");
            });

            options.BaseName.text(BaseData.ProductName);
            options.BasePrice.text($(this).formatfloat(BaseData.Price, 1));
            options.BaseQRCodeLink.attr("href", BaseData.QrCode);
            options.BaseQRCodeImage.attr("src", BaseData.QrCode);
            $("#short").val(BaseData.ShortUrl);
            options.BaseShortUrl.text(BaseData.ShortUrl);

            let specList = "";
            if (BaseData.IsSpec === 1) {
                specList += "<select id='spec'>";
                specList += "<option>請選擇規格</option>"
            } else {
                specList += "<select id='spec' style='display: none'>";
            }
            $.each(BaseData.ProductSpecList, function (index, value) {
                specList += "<option data-id='" + value.ProductSpecId + "' data-quan='" + value.Quantity + "' data-price='" + value.Price + "' >" + value.Spec + "</option>";
            });
            specList += "</select>";
            options.BaseSpecList.append(specList);

            $.each(BaseData.ShippingList, function (k, v) {
                let text;
                if (v.Type.substring(0, 3) === "Cvs") {
                    text = "超商取貨"
                } else {
                    text = v.Text;
                }
                options.BaseShipList.append("<option data-id='" + v.Type + "' data-fee = '" + v.Price + "'>" + text + " NT$ " + v.Price + "</option>");
            });

            let price = options.BaseSpecList.find("option:selected").data("price");
            let Shipping = options.BaseShipList.find("option:selected").data("fee");
            let quantity = options.BaseCountnum.html();
            options.BaseTotal.text($(this).formatfloat(price * quantity + Shipping, 0));

            $("#touchSlider").touchSlider({controls:false});

        };

        plugin.CalculateTotal = function (){
            let price = options.BaseData.Price;
            let Shipping = options.BaseShipList.find("option:selected").data("fee");
            let quantity = options.BaseCountnum.html();
            options.BaseTotal.text($(this).formatfloat(price * quantity + Shipping, 0));
        }

        plugin.SpinnerDownBtn = function() {
            let count = options.BaseCountnum.html();
            if (parseInt(count) === 1) {
                count = 1;
                options.BaseDownBtn.css('background-image', 'url(/static/img/icon-close-cut-12-gray-70.svg)');
            } else {
                count = parseInt(count) - 1;
                options.BaseUpBtn.css('background-image', 'url(/static/img/icon-open-add-12-black-30.svg)');
            }
            options.BaseCountnum.text(count).trigger('change');
        };

        plugin.SpinnerUpBtn = function() {
            let count = options.BaseCountnum.html();
            let quan = options.BaseSpecList.find("option:selected").data("quan");
            if (parseInt(count) === quan) {
                count = quan;
                options.BaseUpBtn.css('background-image', 'url(/static/img/icon-open-add-12-gray-70.svg)');
            } else {
                count = parseInt(count) + 1;
                options.BaseDownBtn.css('background-image', 'url(/static/img/icon-close-cut-12-black-30.svg)');
            }
            options.BaseCountnum.text(count).trigger('change');
        }

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