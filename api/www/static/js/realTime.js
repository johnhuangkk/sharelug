$(function() {
    $.fn.realtime = function (elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseProductName:  $('input[name="name"]'),
            BaseProductPrice: $('input[name="price"]'),
            BaseShippingList: $('select[name="shipping"]'),
            BaseProductShipFee: $('input[name="shipfee"]'),
            BaseProductPayWay: $('select[name="payWay"]'),
            BaseTotal: $('#total'),
            BaseSubmitBtn: $("#submit"),
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        plugin.init = function () {
            options.BaseShippingList.bind("change", function () {
                let elem = $("#shipfee");
                if ($(this).val() === "none") {
                    elem.find("input").remove();
                    elem.append("<input type='text' name='shipfee' class='fee validate' placeholder='請說明交易內容' />");
                } else if ($(this).val() === "F2F") {
                    elem.find("input").remove();
                    elem.append("<input type='text' name='shipfee' class='fee' value='NT$ 0' disabled />");
                } else {
                    elem.find("input").remove();
                    elem.append('<input type="tal" name="shipfee" class="fee validate" maxlength="3" placeholder="NT$" oninput="value=value.replace(/[^\\d]/g,\'\')" />');
                }
                $('input[name="shipfee"]').on("change", function () {
                    options.BaseProductPayWay.focus();
                    plugin.CalculateTotal();
                    plugin.ButtonCheck();
                });
                plugin.ButtonCheck();
            });

            options.BaseProductName.on("change", function (){
                options.BaseProductPrice.focus();
                plugin.ButtonCheck();
            });


            options.BaseProductPayWay.on("change", function () {
                plugin.ButtonCheck();
            });

            options.BaseProductPrice.bind("change", function () {
                let elem = options.BaseShippingList;
                if (elem.children("option:selected").val() === ""){
                    elem.focus();
                }
                plugin.CalculateTotal();
                plugin.ButtonCheck();
            });

            options.BaseShippingList.on("change", function (){
                options.BaseProductShipFee.focus();
                plugin.CalculateTotal();
                plugin.ButtonCheck();
            });

            $(document).on("change", ".fee", function () {
                options.BaseProductPayWay.focus();
                plugin.CalculateTotal();
                plugin.ButtonCheck();
            });

            options.BaseProductPayWay.on("change", function (){
                plugin.ButtonCheck();
            });

            $("#file-upload").on('change', function () {
                let img = this.files[0];
                let reader = new FileReader();
                reader.onloadend = function() {
                    let images = reader.result;
                    $('input[name="image"]').val(encodeURIComponent(images));
                    $('.uploadImgS').find('label').css("background-image", "url("+images+")");
                }
                reader.readAsDataURL(img);
                plugin.ButtonCheck();
            });

            $("#ProductList").bind("click", function (){
                modals("/v1/store/realtime?page=0",
                    {
                        backdrop:true,
                        response: function (e) {
                            $.http.prerequest(function () {
                            }).get("/v1/store/edit/product/"+ $(e).data("id")).done(function (result) {
                                if (result.Status === "Success") {
                                    console.log(result);
                                    plugin.setProduct(result.Data);
                                } else {
                                    console.log(result);
                                }
                            });
                        },
                    }
                );
            });

            options.BaseSubmitBtn.bind("click", function() {

                let images = [];
                if ($('input[name="image"]').val() !== "") {
                    images.push($('input[name="image"]').val());
                }

                console.log($('input[name="shipfee"]').val());

                let shipping = [];
                let ship = {
                    ShipType:  options.BaseShippingList.find(':selected').val(),
                    ShipFee: parseInt($('input[name="shipfee"]').val()),
                    ShipRemark: options.BaseShippingList.find(':selected').val() === "none" ? $('input[name="shipfee"]').val():"",
                }
                shipping.push(ship);

                let payWay = [];
                payWay.push(options.BaseProductPayWay.find(':selected').val());

                let data = {
                    ProductImage: images,
                    ProductName: options.BaseProductName.val(),
                    ProductSpecList: [],
                    ProductQty: 1,
                    IsSpec: 0,
                    Price: options.BaseProductPrice.val(),
                    ShippingList: shipping,
                    PayWayList: payWay,
                    ShipMerge: 0,
                    IsRealtime: 1
                }

                $.http.prerequest(function () {
                }).post("/v1/store/product/new/post", JSON.stringify(data)).done(function (result) {
                    if (result.Status === "Success") {
                        console.log(result.Data)
                        window.location.href = "/store/realtime/review/" + result.Data.ProductId;
                    } else {
                        console.log(result.Message);
                        plugin.alert("錯誤訊息：", result.Message);
                    }
                }).fail(function (result) {
                    console.log(result);
                });

            });

            $("#footer").addClass("hideM");

        };

        plugin.view = function () {

        };

        plugin.alert = function(title, text) {
            let html = '<section class="popupZone"><div class="wrap"><div class="alertIconText">' +
                '<p class="iconAlert"></p><h4>'+title+'</h4><div class="cont"><p>' + text + '</p></div>' +
                '<div class="btns"><span><button class="genBtn h37 prime">確認</button></span>' +
                '</div></div></div></section>';
            $('body').append(html);
            $('.popupZone button').on("click", function () {
                $('body').find('.popupZone').remove();
            });
        };

        plugin.ButtonCheck = function() {
            let buttonChecker = true;
            $("form .validate").each(function (index, elem) {
                if ($(elem).val() === "") {
                    buttonChecker = false;
                }
            });
            if (buttonChecker === true) {
                $("#submit").removeClass("disable").addClass("prime").attr('disabled', false);
            } else {
                $("#submit").removeClass("prime").addClass("disable").attr('disabled', true);
            }
        }

        plugin.CalculateTotal = function (){
            let price = 0;
            let Shipping = 0;
            if (options.BaseProductPrice.val() !== "") {
                price = parseInt(options.BaseProductPrice.val());
            }
            if ($('input[name="shipfee"]').val() !== "") {
                var str = $('input[name="shipfee"]').val();
                var res = str.replace(/NT$/g, "");
                Shipping = res;
                if (options.BaseShippingList.find(':selected').val() === "none") {
                    Shipping = 0;
                }
            }
            let total = parseInt(price) + parseInt(Shipping);
            if (total > 0) {
                options.BaseTotal.removeClass("zero");
            } else {
                options.BaseTotal.addClass("zero");
            }
            options.BaseTotal.text($(this).formatfloat(total, 0));
        };

        plugin.setProduct = function (data) {
            options.BaseProductName.val(data.ProductName);
            options.BaseProductPrice.val(data.Price).trigger('change');
            $.each(data.ShippingList , function (i, v){
                options.BaseShippingList.val(v.Type).trigger('change');

                let elem = $("#shipfee");
                if (v.Type === "none") {
                    elem.find("input").remove();
                    elem.append("<input type='text' name='shipfee' class='fee validate' value='"+v.Remark+"' maxlength='10' placeholder='請說明交易內容' />");
                } else if (v.Type === "F2F") {
                    elem.find("input").remove();
                    elem.append("<input type='text' name='shipfee' class='fee' value='NT$ 0' disabled />");
                } else {
                    elem.find("input").remove();
                    elem.append('<input type="tal" name="shipfee" class="fee validate" maxlength="3" value="'+v.Price+'" placeholder="NT$" oninput="value=value.replace(/[^\\d]/g,\'\')" />');

                }
            });

            $.each(data.ProductPayWayList, function (i, v){
                options.BaseProductPayWay.val(v.Type);
            });

            if (data.ProductImageList !== null) {
                $('input[name="image"]').val(data.ProductImageList[0].split('/')[6]);
                $('.uploadImgS').find('label').css("background-image", "url("+data.ProductImageList[0]+")");
            } else {
                $('input[name="image"]').val("");
                $('.uploadImgS').find('label').css("background-image", "url('/static/img/upLoadImgBg.svg')");
            }
            plugin.CalculateTotal();
            plugin.ButtonCheck();
        };

        plugin.init();
    }
});