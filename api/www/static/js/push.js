$(function () {
    $.fn.push = function (opt) {
        let DEFAULT_OPTIONS = {
            BaseUploadBtn: $("#file-upload"),
            BaseImageFile: $("input[name='image']"),
            BaseSubmitBtn: $("#submit"),
            BaseSpecAddBtn: $("#specAdd"),
            BaseAddSpecs: $(".addSpec"),
            BaseUploadItems: $("#items"),
            BaseSpecList: $(".formS02"),
            BaseProductName:  $('input[name="name"]'),
            BaseProductAmt: $('input[name="price"]'),
            BaseProductQty: $('input[name="qty"]'),
            BaseProductShip: $("#ship .items"),
            BaseProductPayWay: $("#payWay .items"),
            BaseSwitchBtn: $(".switch"),
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        let max_spec = 5;

        plugin.init = function () {
            plugin.reload();

            options.BaseSwitchBtn.on("click", function (){
                // console.log($(this).attr("class"));
                $(this).parent().parent().find("div").slideToggle(function () {
                    console.log($(this).attr("class"));
                });
            });

            options.BaseProductName.on("click", function(){
                $(this).removeClass('errorNoImg').parent().find(".errorAlert").remove();
            });

            options.BaseProductAmt.on("click", function () {
                $(this).removeClass('errorNoImg').parent().find(".errorAlert").remove();
            });

            options.BaseProductQty.on("click", function () {
                $(this).removeClass('errorNoImg').parent().find(".errorAlert").remove();
            });

            options.BaseSpecList.on("click", ".addSpec", function () {
                $('.qty').hide();
                options.BaseSpecList.find("a").remove();
                let html = "<h5>請填寫規格及數量，最多 5 組：</h5>" +
                    "<p class='items'>\n" +
                    "<span class='specTitle'><input type='text' class='spec' name='productspec' placeholder='規格（限 14 字）' maxlength='14'></span>\n" +
                    "<span class='specQty'><input type='text' class='qty' name='quantity' placeholder='1~99' oninput=\"value=value.replace(/[^\\d]/g,'')\" maxlength='2'></span>\n" +
                    "<span class='del' title='刪除'></span>\n" +
                    "</p>\n" +
                    "<a class='genBtnBorder' id='specAdd' href='javascript:void(0)'>新增規格及數量</a>";
                options.BaseSpecList.append(html);
            });

            options.BaseUploadItems.on("click", ".close", function(element) {
                element.preventDefault();
                $(this).parent().closest("div").remove();
                plugin.reload();
            });

            options.BaseSpecList.on("click", "#specAdd", function(element) {
                element.preventDefault();
                if (plugin.checkSpec()) {
                    return false;
                }
                let spec_count = $('.formS02 .items').length;
                if(spec_count < max_spec) {
                    let html = "<p class='items'>\n" +
                        "<span class='specTitle'><input type='text' class='spec' name='productspec' placeholder='規格（限 14 字）' maxlength='14'></span>\n" +
                        "<span class='specQty'><input type='text' class='qty' name='quantity' placeholder='1~99' oninput=\"value=value.replace(/[^\\d]/g,'')\" maxlength='2'></span>\n" +
                        "<span class='del' title='刪除'></span>\n" +
                        "</p>";
                    $(this).before(html);
                }
            });

            options.BaseSpecList.on("click", ".del", function(element) {
                element.preventDefault();
                $(this).parents("p").remove();
                let spec_count = $('.formS02 .items').length;
                if (spec_count === 0) {
                    $('.formS02').empty().append("<a class='addSpec' href='javascript:void(0)'>新增規格</a>");
                    $('.formS01 .qty').show();
                }
            });

            options.BaseSpecList.on("click", ".spec", function () {
                $(this).removeClass("errorNoImg");
                $(this).parent().next().find(".qty").removeClass("errorNoImg");
                $(this).parent().parent().find(".errorAlert").remove();
            })

            options.BaseSpecList.on("click", ".qty", function () {
                $(this).removeClass("errorNoImg");
                $(this).parent().prev().find(".spec").removeClass("errorNoImg");
                $(this).parent().parent().find(".errorAlert").remove();
            });

            $("#ProductList").bind("click", function (){
                modals("/v1/store/products?length=10&page=1",
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

            $('#shipDesc').on("click", function () {
                let html = '<p>1.不同配送方式，不能放入同一張結帳清單中。</p>' +
                    '<p>2.結帳清單中，分開列出「合併運費」及「不合併運費」的商品。合併運費的商品以 運費最高為主。</p><h5>舉例一</h5>' +
                    '<p>可合併運費商品 A 運費 NT$ 65 與可合併運費商品 B 運費 NT$ 80，此時結帳清單的運費會以 NT$ 80 為計價。</p><h5>舉例二</h5>' +
                    '<p>可合併運費商品 A 運費 NT$ 65 與不可合併運費商品 H 運費 NT$ 80，此時結帳清單的運費會以 NT$ 145 為計價。</p><h5>舉例三</h5>' +
                    '<p>不可合併運費商品 S 運費 NT$ 65 與不可合併運費商品 Q 運費 NT$ 90，此時結帳清單的運費會以 NT$ 155 為計價。</p>'
                alerts("運費說明", html);
            });

            $('.iconInfo').on("click", function () {
                let html = '<p>1.不同配送方式，不能放入同一張結帳清單中。</p>' +
                    '<p>2.結帳清單中，分開列出「合併運費」及「不合併運費」的商品。合併運費的商品以 運費最高為主。</p><h5>舉例一</h5>' +
                    '<p>可合併運費商品 A 運費 NT$ 65 與可合併運費商品 B 運費 NT$ 80，此時結帳清單的運費會以 NT$ 80 為計價。</p><h5>舉例二</h5>' +
                    '<p>可合併運費商品 A 運費 NT$ 65 與不可合併運費商品 H 運費 NT$ 80，此時結帳清單的運費會以 NT$ 145 為計價。</p><h5>舉例三</h5>' +
                    '<p>不可合併運費商品 S 運費 NT$ 65 與不可合併運費商品 Q 運費 NT$ 90，此時結帳清單的運費會以 NT$ 155 為計價。</p>'
                alerts("運費說明", html);
            });

            options.BaseSubmitBtn.bind("click", function(){
                let validate = false;
                if (options.BaseProductName.val() === ""){
                    options.BaseProductName.parent().find(".errorAlert").remove();
                    options.BaseProductName.addClass("errorNoImg");
                    options.BaseProductName.parent().append("<div class=\"errorAlert\">請填寫商品名稱。</div>");
                    options.BaseProductName.focus();
                    validate = true;
                }

                if ($(this).validateChangeCount(options.BaseProductName.val()) > 40) {
                    options.BaseProductName.parent().find(".errorAlert").remove();
                    options.BaseProductName.addClass("errorNoImg");
                    options.BaseProductName.parent().append("<div class=\"errorAlert\">商品名稱不可大於40個字。</div>");
                    options.BaseProductName.focus();
                    validate = true;
                }

                if (options.BaseProductAmt.val() === "") {
                    options.BaseProductAmt.parent().find(".errorAlert").remove();
                    options.BaseProductAmt.addClass("errorNoImg");
                    options.BaseProductAmt.parent().append("<div class=\"errorAlert\">請填寫商品金額。</div>");
                    options.BaseProductAmt.focus();
                    validate = true;
                }

                if (!options.BaseProductQty.is(":hidden") && options.BaseProductQty.val() === "") {
                    options.BaseProductQty.parent().find(".errorAlert").remove();
                    options.BaseProductQty.addClass("errorNoImg");
                    options.BaseProductQty.parent().append("<div class=\"errorAlert\">請填寫商品數量。</div>");
                    options.BaseProductQty.focus();
                    validate = true;
                }

                if (options.BaseProductQty.is(":hidden")) {
                    if (plugin.checkSpec()) {
                        validate = true;
                    }
                }
                options.BaseProductShip.parent().find(".errorAlert").remove();
                $.each(options.BaseProductShip, function (index, element) {

                    if ($(element).find(".shipType").is(":checked") && $(element).find(".shipFee").attr("disabled") === undefined) {
                        if ($(element).find(".shipFee").val() === "") {
                            $(element).find(".spec").addClass("errorNoImg");
                            $(element).next().find(".qty").addClass("errorNoImg");
                            $(element).parent().append("<div class=\"errorAlert\">請填寫運費金額。</div>");
                            $(element).focus();
                            validate = true;
                        }
                    }
                });

                let shipping = [];
                $.each(options.BaseProductShip, function (index, element) {
                    let shipType = $(element).find(".shipType");
                    let shipFee = $(element).find(".shipFee");
                    if (shipType.is(":checked")){
                        let ship = {
                            ShipType:  shipType.val(),
                            ShipFee: parseInt(shipFee.val()),
                        }
                        shipping.push(ship);
                    }
                });

                console.log(shipping, shipping.length);

                let payWay = [];
                $.each(options.BaseProductPayWay, function (index, element) {
                    let type = $(element).find(".type");
                    if (type.is(":checked")) {
                        payWay.push(type.val());
                    }
                });

                let images = [];
                $.each($("input[name='image']"), function (index, element) {
                    images.push($(element).val());
                });

                if (shipping.length === 0) {
                    plugin.alert("錯誤訊息：", "請確認配送方式及運費");
                    validate = true;
                }

                if (payWay.length === 0) {
                    plugin.alert("錯誤訊息：", "請確認付款方式");
                    validate = true;
                }

                if (images.length === 0) {
                    plugin.alert("錯誤訊息：", "請至少選擇一張圖片");
                    validate = true;
                }

                let ProductSpecList = [];
                if (options.BaseProductQty.is(":hidden")) {
                    plugin.checkSpec();
                    $.each($('.formS02 .items'), function (index, element) {
                        let spec = {
                            ProductSpec: $(element).find(".spec").val(),
                            Quantity: parseInt($(element).find(".qty").val()),
                        }
                        ProductSpecList.push(spec);
                    })
                }

                if (!validate) {
                    //產生data
                    let data = {
                        ProductImage: images,
                        ProductName: options.BaseProductName.val(),
                        ProductSpecList: ProductSpecList,
                        ProductQty: parseInt(options.BaseProductQty.val()),
                        IsSpec: options.BaseProductQty.is(":hidden") ? 1:0,
                        Price: options.BaseProductAmt.val(),
                        ShippingList: shipping,
                        PayWayList: payWay,
                        ShipMerge: $("input[name='shipMerge']").prop("checked") ? 1:0
                    }
                    // console.log(data);
                    $.http.prerequest(function () {
                    }).post("/v1/store/product/new/post", JSON.stringify(data)).done(function (result) {
                        if (result.Status === "Success") {
                            console.log(result.Data)
                            window.location.href = "/store/product/review/" + result.Data.ProductId;
                        } else {
                            console.log(result.Message);
                            plugin.alert("錯誤訊息：", result.Message);
                        }
                    }).fail(function (result) {
                        console.log(result);
                    });
                }
            });
        };

        plugin.checkSpec = function() {
            let err = false;
            $.each($(".specTitle"), function (index, element) {
                if ($(element).find(".spec").val() === "" || $(element).next().find(".qty").val() === "" || $(element).next().find(".qty").val() < 1) {
                    $(element).parent().find(".errorAlert").remove();
                    $(element).find(".spec").addClass("errorNoImg");
                    $(element).next().find(".qty").addClass("errorNoImg");
                    $(element).parent().append("<div class=\"errorAlert\">請填寫規格及數量。</div>");
                    err = true;
                }
            });
            return err
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

        plugin.reload = function(image = ""){
            let html = "";
            let count = $(".item").length;
            if (count > 0) {
                $.each($(".item"), function (index, element) {
                    html = html + $(element)[0].outerHTML;
                })
            }
            options.BaseUploadItems.find(".box").remove();
            if (image !== "") {
                html = html + image;
                count += 1;
            }
            if (count < 8 ) {
                html = html + '<div class="box ignore"><span class="addPic">' +
                    '<label for="file-upload" class="custom-file-upload">' +
                    '<input id="file-upload" type="file" accept="image/*"/>新增照片' +
                    '</label>' +
                    '</span></div>';
            }
            for (i=0;i<=(6 - count); i++) {
                html = html + "<div class='box ignore'><span>"+(i + 2 + count)+"/8</span></div>";
            }
            options.BaseUploadItems.append(html);
            $("#file-upload").on('change', function () {
                let img = this.files[0];
                let reader = new FileReader();
                reader.onloadend = function() {
                    let images = reader.result;
                    let image = "<div class='box item'>" +
                        "<span class='close'></span>" +
                        "<span class='pic'><img src='" + images + "'/>" +
                        "<input type='hidden' name='image' class='image' value='" + encodeURIComponent(images)+ "' /></span></div>";
                    plugin.reload(image);
                }
                reader.readAsDataURL(img);
            });
        };

        plugin.setProduct = function(data) {
            $("#items").find("div").remove();
            $.each(data.ProductImageList, function (i, v){
                let image = "<div class='box item'>" +
                    "<span class='close'></span>" +
                    "<span class='pic'><img src='" + v + "'/>" +
                    "<input type='hidden' name='image' class='image' value='" + v.split('/')[6] + "' /></span></div>";
                plugin.reload(image);
            });

            options.BaseProductName.val(data.ProductName);
            options.BaseProductAmt.val(data.Price);

            if (data.IsSpec === 1) {
                $('.qty').hide();
                options.BaseSpecList.find("a").remove();
                options.BaseSpecList.empty();
                let html = "<h5>請填寫規格及數量，最多 5 組：</h5>";
                $.each(data.ProductSpecList, function (i, v) {
                    html += "<p class='items'>\n" +
                        "<span class='specTitle'><input type='text' class='spec' name='productspec' value='"+v.Spec+"' placeholder='規格（限 14 字）' maxlength='14'></span>\n" +
                        "<span class='specQty'><input type='text' class='qty' name='quantity' value='1' placeholder='1~99' oninput=\"value=value.replace(/[^\\d]/g,'')\" maxlength='2'></span>\n" +
                        "<span class='del' title='刪除'></span>\n" +
                        "</p>\n";
                });

                html += "<a class='genBtnBorder' id='specAdd' href='javascript:void(0)'>新增規格及數量</a>";
                options.BaseSpecList.append(html);
            } else {
                options.BaseProductQty.val(1);
            }
        };

        plugin.init();
    }
});