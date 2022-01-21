$(function(){
    $.fn.products = function(elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseImages: $("#product_image"),
            BaseImageBtn: $("#product_image_btn"),
            BaseName: $("#product_name"),
            BaseSpecList: $(".spec .select"),
            BaseShipping: $("#shipping"),
            BaseShipList: $("#shippinglist"),
            BasePaywayList: $("#paywaylist"),
            BaseDownBtn: $(".down"),
            BaseUpBtn: $(".up"),
            BaseCountnum: $("#countnum"),
            BasePrice: $("#price"),
            BasePayBtn: $("#pay"),
            BaseQrcode: $("#qrcode"),
            BaseTotal: $("#total"),
            BaseData:"",
            BaseCartCount: 0,
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        plugin.init = function() {
            plugin.view(options.BaseData);

            options.BaseDownBtn.click(function () {
                plugin.SpinnerDownBtn();
            });
            options.BaseUpBtn.click(function () {
                plugin.SpinnerUpBtn();
            });

            options.BaseCountnum.bind("change", function () {
                plugin.CalculateTotal();
            });

            options.BaseShipping.bind("change", function () {
                plugin.CalculateTotal();
            });

            options.BaseSpecList.bind("change", function () {
                plugin.specChange(this);
            });
            options.BasePayBtn.click(function () {
                plugin.addCartBtn();
            });

            $(document).on("click", ".hCart", function () {
                location.href = "/pay";
            });
        };

        plugin.view = function(BaseData) {
            options.BaseName.text(BaseData.ProductName);
            $.each(BaseData.ProductImageList, function (k, v) {
                options.BaseImages.append("<li style=\"background-image: url(" + v + ")\"></li>");
            });

            options.BaseQrcode.attr("src", '/static/images/qrcode/'+BaseData.ProductId+'.jpg');
            let qty = 0;
            let specList = "";
            if (BaseData.IsSpec === 1) {
                specList += "<select id='spec'>";
                specList += "<option>請選擇規格</option>"
            } else {
                specList += "<select id='spec' style='display: none'>";
            }
            $.each(BaseData.ProductSpecList, function (index, value) {
                if (value.Quantity > 0) {
                    specList += "<option data-id='" + value.ProductSpecId + "' data-qty='" + value.Quantity + "' data-price='" + value.Price + "' >" + value.Spec + "</option>";
                    qty = value.Quantity;
                }
            });
            specList += "</select>";
            options.BaseSpecList.append(specList);
            console.log(BaseData);
            let c = "";
            $.each(BaseData.ShippingList, function (k, v) {
                if (v.Type.substring(0, 3) === "Cvs") {
                    c = "超商取件";
                } else {
                    c = "";
                }
                options.BaseShipping.append("<option data-id='" + v.Type + "' data-fee = '" + v.Price + "'>" + v.Text + c + " NT$ " + v.Price + "</option>");
            });
            options.BasePrice.text($(this).formatfloat(BaseData.Price, 1));

            options.BaseShipList.val(BaseData.ShippingList);
            options.BasePaywayList.val(BaseData.PaywayList);

            let price = options.BaseData.Price;
            let Shipping = options.BaseShipping.find("option:selected").data("fee");
            let quantity = options.BaseCountnum.html();
            options.BaseTotal.text($(this).formatfloat(price * quantity + Shipping, 0));
            plugin.getCartsCount(qty);


            $("#touchSlider").touchSlider({controls:false});
        };

        plugin.CalculateTotal = function (){
            let price = options.BaseData.Price;
            let Shipping = options.BaseShipping.find("option:selected").data("fee");
            let quantity = options.BaseCountnum.html();
            options.BaseTotal.text($(this).formatfloat(price * quantity + Shipping, 0));
        }

        plugin.SpinnerDownBtn = function() {
            let count = options.BaseCountnum.html();
            if (parseInt(count) === 1) {
                count = 1;
                options.BaseDownBtn.css('background-image', 'url(/static/img/icon-close-cut-12-gray-70.svg)');
            } else if(count > 1) {
                count = parseInt(count) - 1;
                options.BaseUpBtn.css('background-image', 'url(/static/img/icon-open-add-12-black-30.svg)');
            }
            options.BaseCountnum.text(count).trigger('change');
        };

        plugin.SpinnerUpBtn = function() {
            let count = options.BaseCountnum.html();
            let qty = options.BaseSpecList.find("option:selected").data("qty");
            if (parseInt(count) === qty) {
                count = qty;
                options.BaseUpBtn.css('background-image', 'url(/static/img/icon-open-add-12-gray-70.svg)');
            } else if (count <= qty) {
                count = parseInt(count) + 1;
                options.BaseDownBtn.css('background-image', 'url(/static/img/icon-close-cut-12-black-30.svg)');
            }
            options.BaseCountnum.text(count).trigger('change');
        }

        plugin.specChange = function() {
            options.BasePrice.text(options.BaseSpecList.find("option:selected").data("price"));
            options.BaseCountnum.text(1).trigger('change');
            options.BaseDownBtn.css('background-image', 'url(/static/img/icon-close-cut-12-gray-70.svg)');
            options.BaseUpBtn.css('background-image', 'url(/static/img/icon-open-add-12-black-30.svg)');
        }

        plugin.addCartBtn = function() {
            let validate = false;

            $('body').find(".errorAlert").remove();
            if (options.BaseSpecList.find("option:selected").data("id") === undefined) {
                options.BaseSpecList.find('select').addClass("errorNoImg");
                options.BaseSpecList.append("<div class=\"errorAlert\">請選擇規格。</div>");
                validate = true;
            }

            if (options.BaseShipping.find("option:selected").data("id") === undefined) {
                options.BaseShipping.addClass("errorNoImg");
                options.BaseShipping.parent().append("<div class=\"errorAlert\">請選擇配送方式。</div>");
                validate = true;
            }

            if (!validate) {
                let data = {
                    ProductSpecId: options.BaseSpecList.find("option:selected").data("id"),
                    Quantity: parseInt(options.BaseCountnum.html()),
                    Shipping: options.BaseShipping.find("option:selected").data("id"),
                }
                plugin.addCart(data);
            }
        }

        plugin.addCart = function(data) {
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
        };

        plugin.getCartsCount = function(qty){
            $.http.prerequest(function () {
            }).get('/v1/cart/count')
                .done(function (result) {
                    console.log("sss", qty);
                    if (qty === 0) {
                        $("#countnum").text(0);
                        $("#pay").removeClass("prime").addClass("disable").html("商品已售完");
                    } else if (result.Data.Count !== 0) {
                        options.BasePayBtn.html("加入結帳清單")
                    } else {
                        options.BasePayBtn.html("立即結帳")
                    }

                    if (result.Data.Count !== 0) {
                        $("#right").append("<button class='hCart full'></button>");
                    } else {
                        $("#right").append("<button class='hCart'></button>");
                    }

                }).fail(function (result) {
                window.location.href="/404";
            });
        };

        plugin.init();
    }
});


