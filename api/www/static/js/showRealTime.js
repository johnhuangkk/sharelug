$(function() {
    $.fn.RealTimeShow = function (elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseProductName: $("#productName"),
            BaseProductPrice: $("#ProductPrice"),
            BaseShipText: $("#ShipText"),
            BaseShipFee: $("#ShipFee"),
            BasePayWayText: $("#PayWayText"),
            BaseTotal: $("#Total"),
            BaseUpLoadImg: $(".uploadImgS"),
            BaseSeller: $(".seller"),
            BaseData:"",
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        plugin.init = function (){
            plugin.view();

            $("#pay").on("click", function () {
                let data = {
                    ProductSpecId: $("input[name='ProductSpecId']").val(),
                    Quantity: 1,
                    Shipping: $("input[name='Shipping']").val(),
                }
                console.log(data);
                $.http.prerequest(function () {
                }).put('/v1/cart/add', JSON.stringify(data)).done(function (result) {
                    if (result.Status === "Success") {
                        window.location.href="/pay";
                    } else {
                        console.log(result);
                        alert(result.Message);
                    }
                }).fail(function (result) {
                    window.location.href="/404";
                });
            });

        };

        plugin.view = function () {
            console.log(options.BaseData);
            let shipfee = 0;

            options.BaseProductName.text(options.BaseData.ProductName);
            options.BaseProductPrice.text(options.BaseData.Price);


            options.BaseSeller.find("a").attr("href", "/store/list/" + options.BaseData.StoreId);
            options.BaseSeller.find("img").attr("src", options.BaseData.StorePicture);
            options.BaseSeller.find("span").text(options.BaseData.StoreName);


            $.each(options.BaseData.ShippingList, function (index, value) {
                options.BaseShipText.text(value.Text);
                options.BaseShipFee.text($(this).formatfloat(value.Price));
                shipfee = value.Price;
            });

            options.BaseTotal.text($(this).formatfloat(options.BaseData.Price + shipfee));

            $.each(options.BaseData.ProductPayWayList, function (index, value) {
                options.BasePayWayText.text(value.Text);
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
                $(".billOffical").addClass("noPic");
            }


            let html = "<input type='hidden' name='ProductSpecId' value='"+options.BaseData.ProductSpecList[0].ProductSpecId+"' />"
                + "<input type='hidden' name='Shipping' value='"+options.BaseData.ShippingList[0].Type+"' />";
            $("#pay").append(html);

            $("#pageHeader").hide();
            $("#footer").addClass("hideM");

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
