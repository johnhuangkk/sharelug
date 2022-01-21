$(function () {
    $.fn.pay = function (elem, opt) {
        let DEFAULT_OPTIONS = {
            BaseStoreList: $("#store"),
            BaseShipList: $("#shipping"),
            BaseCardInfo: $("#cardInfo"),
            BaseShippingBtn: $(".btn-ship"),
            BaseSubtotal: $("#subtotal"),
            BaseShipFee: $("#shipfee"),
            BaseTotal: $("#total"),
            BaseShipText: $("#ship_text"),
            BasePayWayList: $("#paywaylist"),
            BaseShipInfo: $("#shipInfo"),
        };

        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);
        let shipfee;
        let payway = "";
        plugin.init = function () {
            $.http.prerequest(function () {
            }).get("/v1/cart/get").done(function (result) {
                if (result.Status === "Success") {
                    // console.log(result.Data);
                    plugin.view(result.Data);
                } else {
                    alerts("", result.Message, {
                        confirm: function() {
                            window.location.href="/404";
                        }
                    });
                }
            }).fail(function (result) {
                console.log(result);
            });

            $("#orderer").on('change',function () {
                let e = $(this);
                $("#extra").slideToggle(function () {
                    if($(this).is(":visible")) {
                        e.prop('checked', false);
                        $(this).find('input').attr('disabled', false);
                    } else {
                        e.prop('checked', true);
                        $(this).find('input').attr('disabled', 'disabled');
                    }
                });
            }).trigger('click');

            $(document).on('click', '.down', function () {
                let quantity = parseInt($(this).next().text()) - 1;
                let id = $(this).data('id');
                if (quantity === 0) {
                    confirms("","確定刪除此商品！！", {
                        confirm: function() {
                            $(this).next().text(quantity);
                            plugin.DeleteProduct(id);
                        }
                    });
                } else if (quantity > 0){
                    $(this).next().text(quantity);
                    plugin.ChangeQuantity("delete", $(this).data('id'), $(this).next(), quantity);
                }
            });

            $(document).on('click', '.up', function () {
                let quantity = parseInt($(this).prev().text()) + 1;
                $(this).prev().text(quantity);
                plugin.ChangeQuantity("add", $(this).data('id'), $(this).prev(), quantity);
            });

            $("#switch").on("click", function(){
                let e = $(this);
                $("#cardInfo").slideToggle(function () {
                    if($(this).is(":visible")){
                        e.addClass("off");
                    }else{
                        e.removeClass("off");
                    }
                });
            }).trigger('click');

            $("#submit").on("click", function (e) {
                e.preventDefault();
                let data = $("form").serializeObject();
                if (checker === false ) {
                    plugin.submitPhone($("input[name='BuyerPhone']").val(), false);
                    return;
                }

                $.each($("input[name='ReceiverAddress']"), function (index, element) {
                    if ($(element).attr("disabled") === undefined) {
                        $(element).removeClass("error").next(".errorAlert").remove();
                        let addr = plugin.getChineseLength(element);
                        if (addr < 4) {
                            $(element).addClass("error").after("<div class='errorAlert'>地址輸入不正確！</div>");
                            return;
                        }
                    }
                })
                if (payway === "credit") {
                    let elem2 = $("input[name='CreditNumber']");
                    re = /^\d{4}\s\d{4}\s\d{4}\s\d{4}$/;
                    if (!re.test(elem2.val())) {
                        elem2.addClass("error").after("<div class='errorAlert'>信用卡輸入不正確！</div>");
                        return;
                    }
                    let elem3 = $("input[name='CreditExpiration']");
                    re = /^\d{2}\/\d{2}$/;
                    if (!re.test(elem3.val())) {
                        elem3.addClass("error").after("<div class='errorAlert'>信用卡有效年月輸入不正確！</div>");
                        return;
                    }

                    let elem4 = $("input[name='CreditSecurity']");
                    if (elem4.val().length < 3) {
                        elem4.addClass("error").after("<div class='errorAlert'>信用卡安全碼輸入不正確！</div>");
                        return;
                    }
                }
                data.PayWay = payway;

                $.http.prerequest(function () {
                }).post("/v1/cart/pay", JSON.stringify(data)).done(function (result) {
                    if (result.Status === "Success") {
                        window.location.href="/pay/succ/"+ result.Data.OrderId;
                    } else {
                        window.location.href="/pay/fail";
                    }
                }).fail(function (result) {
                    console.log(result);
                });
            });

            // console.log($(".clums > select:last-child").attr("name"));
        };

        plugin.getChineseLength = function(element) {
            let tmp = $(element).val();
            tmp = tmp.replace(/[A-Za-z0-9]/g,"");
            tmp = tmp.replace(/[^\u0000-\u00ff]/g," a ");
            tmp = tmp.replace(/\s/g,"");
            return tmp.length;
        }

        plugin.view = function(BaseData) {

            if (BaseData.Product === null) {
                alerts("","購物車內無商品資料！！",{
                    confirm: function (){
                        history.go(-1);
                    }
                });
            }

            console.log(BaseData.Product[0].Quantity);

            if (BaseData.Product.length === 1 && BaseData.Product[0].Quantity === 1) {
                 options.BaseShipInfo.hide();
                 $("#csvstore").hide();
            }

            plugin.setShipList(BaseData, BaseData.Shipping);
            plugin.setPayway(BaseData.PaywayList);
            plugin.getProduct(BaseData);

            options.BaseSubtotal.text($(this).formatfloat(BaseData.Subtotal, 0));
            options.BaseShipFee.text($(this).formatfloat(BaseData.Shipfee, 0));
            options.BaseTotal.text($(this).formatfloat(BaseData.Total, 0));

            $("#submit span").text($(this).formatfloat(BaseData.Total, 0));

            $(".ship li").on("click", function () {
                // $.Alerts.Confirm("重新設定配送方式為「'+$(this).text()+'」？", "111", "確定", alert("OK"));

                // $.confirm({
                //     title: "",
                //     content: '重新設定配送方式為「'+$(this).text()+'」？',
                //     buttons: {
                //         "確定": function () {
                //             plugin.ChangeShipping($(this).data("type"))
                //         },
                //         "取消": function () {
                //         }
                //     }
                // });
                plugin.ChangeShipping($(this).data("type"))

            });

            $("#paywaylist li").on("click", function () {
                $("#paywaylist li").removeClass("actived");
                $(this).addClass("actived");
                let choose = $(this).data('type');
                payway = choose;
                $(".payway").hide().find('input').attr('disabled', true);
                $(document).find("[data-payway='"+choose+"']").show().find("input").attr('disabled', false);;
                plugin.ButtonCheck();
            });

            payway = payway !== "" ? payway : BaseData.PaywayList[0].Type;

            $("#paywaylist").find("[data-type='"+payway+"']").trigger('click');

            plugin.formCheck();
        };

        plugin.ChangeShipping = function (shipType) {
            let data = {ShipType:shipType};
            $.http.prerequest(function () {
            }).put("/v1/cart/ship/change", JSON.stringify(data)).done(function (result) {
                if (result.Status === "Success") {
                    // console.log(result.Data)
                    plugin.view(result.Data);
                } else {
                    alerts("", result.Message, {
                        confirm:function() {
                            window.location.href="/404";
                        }
                    });
                }
            }).fail(function (result) {
                console.log(result);
            });
            // $("#submit").removeClass("prime").addClass("disable").attr('disabled', true);
        };

        plugin.ChangeQuantity = function(type, productId, event, quantity) {
            let data = {Type: type, ProductSpecId:productId};
            $.http.prerequest(function () {
            }).put("/v1/cart/quantity/change", JSON.stringify(data)).done(function (result) {
                if (result.Status === "Success") {
                    plugin.view(result.Data);
                } else {
                    if (type === "add") {
                        event.text(quantity - 1);
                    }
                    alert(result.Message);
                }
            }).fail(function (result) {
                console.log(result.responseText);
            });
        };

        plugin.DeleteProduct = function(productId) {
            let data = {ProductSpecId:productId};
            $.http.prerequest(function () {
            }).put("/v1/cart/delete", JSON.stringify(data)).done(function (result) {
                if (result.Status === "Success") {
                    plugin.view(result.Data);
                } else {
                    alerts("", result.Message, {
                        confirm:function() {
                            history.go(-1);
                        }
                    });
                }
            }).fail(function (result) {
                console.log(result.responseText);
            });
        }

        let checker = false;
        plugin.submitPhone = function(phone, login) {
            let data = {Phone:phone};
            $.http.prerequest(function () {
            }).put("/v1/cart/otp", JSON.stringify(data)).done(function (result) {
                if (result.Status === "Success") {
                    if (result.Data.Otp === 1) {
                        plugin.popup(phone, login);
                    } else {
                        checker = true;
                    }
                } else {
                    alert(result.Message);
                }
            }).fail(function (result) {
                console.log(result);
            });
        }

        plugin.cardTotal = function(){
            let subtotal = 0;
            $.each(options.BaseCardInfo.find(".row"), function () {
                subtotal = subtotal + $(this).data('price') * $(this).find('#countnum').text();
            });
            options.BaseSubtotal.text("NT$" + $(this).formatfloat(subtotal, 0));
            options.BaseShipFee.text("NT$" + $(this).formatfloat(shipfee, 0));
            options.BaseTotal.text("NT$" + (subtotal + $(this).formatfloat(shipfee, 0)));
        };

        plugin.getProduct = function(BaseData){
            options.BaseCardInfo.html("");
            let m = "";
            let n = "";

            $.each(BaseData.Product, function(k, data){
                if (data.ShipMerge !== 1) {
                    n += "<div class='receiptWrap'><ul class='list'><li data-id='"+data.ProductSpecId+"' data-price='"+data.Price+"'>" +
                        "<p class='title'><a href='#'>"+data.ProductName+"</a></p>\n" +
                        "<p class='spec'>"+data.ProductSpec+"</p>" +
                        "<p class='fAmount'>" +
                        "<span data-id='"+data.ProductSpecId+"' class='down'></span>" +
                        "<span class='min2'>" + data.Quantity + "</span>" +
                        "<span data-id='"+data.ProductSpecId+"' class='up'></span></p>" +
                        "<p class='pPrice'><span class='price tw'>"+ $(this).formatfloat(data.Price, 0) +"</span></p></li></ul>" +
                        "<div class='itemShipping'><p class='title'>合併運費<span></span></p>" +
                        "<p class='pPrice'><span class='price tw'>"+data.ShipFee+"</span></p></div></div>";
                } else {
                    m += "<li data-id='"+data.ProductSpecId+"' data-price='"+data.Price+"'>" +
                        "<p class='title'><a href='#'>"+data.ProductName+"</a></p>" +
                        "<p class='spec'>"+data.ProductSpec+"</p>" +
                        "<p class='fAmount'>" +
                        "<span data-id='"+data.ProductSpecId+"' class='down'></span>" +
                        "<span class='min2'>" + data.Quantity + "</span>" +
                        "<span data-id='"+data.ProductSpecId+"' class='up'></span></p>" +
                        "<p class='pPrice'><span class='price tw'>"+ $(this).formatfloat(data.Price, 0) +"</span></p>" +
                        "</li>";
                }
            })
            if (m !== "") {
                 m = "<div class='receiptWrap'><ul class='list'>"+ m +"</ul></div>";
            }
            options.BaseCardInfo.append(m + n);
        };


        plugin.setShipList = function(BaseData, choose) {
            options.BaseShipList.html("");
            options.BaseStoreList.html("");

            ShipList = BaseData.ShipList.sort(function (a, b) {
                return a.Type > b.Type ? 1 : -1;
            });
            let s = 0;
            if (choose === "store") {
                choose = "Cvs7Eleven";
            }
            $.each(BaseData.ShipList, function (k, v) {
                let c = "";
                if ( choose === v.Type ) {
                    c = "class='actived'";
                    let text;
                    if (v.Type.substring(0, 3) === "Cvs") {
                        text = "超商取貨"
                    } else {
                        text = v.Text;
                    }
                    options.BaseShipText.html( text+"收件資料" + " <a class='back2Shop'  href='/store/list/"+BaseData.StoreId+"' >繼續購物</a>");
                }

                let c1 = "";
                if (choose.substring(0, 3) === "Cvs") {
                    c1 = "class='actived'";
                }

                if (v.Type.substring(0, 3) === "Cvs") {
                    options.BaseStoreList.append("<li data-type='" + v.Type + "' data-fee='" + v.Price + "' " + c + ">" + v.Text + "</li>")
                    if (s === 0) {
                        options.BaseShipList.append("<li data-type='store' " + c1 + ">超商取貨</li>");
                        s++;
                    }
                } else {
                    options.BaseShipList.append("<li data-type='" + v.Type + "' data-fee='" + v.Price + "' " + c + ">" + v.Text + "</li>")
                }
            });


            $(".ship-address").each(function (index, element) {
                let type = ['CvsFamily', 'Cvs7Eleven', 'CvsFamily', 'CvsHiLife', 'CvsOkMart'];
                let choosetype;
                if($.inArray(choose, type) === -1) {
                    choosetype = choose;
                    $("#csvstore").hide();
                } else {
                    choosetype = "Store";
                    if (BaseData.Product.length !== 1) {
                        $("#csvstore").show();
                    }
                }
                if ($(element).data("ship") === choosetype) {
                    $(element).show("slow", function () {
                        $(element).find(".twzipcode").twzipcode({
                            zipcodeIntoDistrict: true,
                            css: ["validate select", "validate select"],
                            countyName: "ReceiverCity", // 自訂城市 select 標籤的 name 值
                            districtName: "ReceiverArea", // 自訂區別 select 標籤的 name 值
                            zipcodeName: "Zipcode",
                        });
                    });
                    $(element).find("input, select").attr('disabled', false);
                } else {
                    $(element).find(".twzipcode").twzipcode('destroy');
                    $(element).hide().find("input, select").attr('disabled', true);
                }
            });

        };

        plugin.setPayway = function(PaywayList){
            options.BasePayWayList.html("");

            let type = ["Credit", "Transfer", "CVSPay", "Balance"];

            $.each(type, function (index, value) {
                console.log(value);
                $.each(PaywayList, function (k, v) {
                    if (v.Type === value) {
                        options.BasePayWayList.append("<li data-type='" + v.Type + "'>" + v.Text + "</li>");
                    }
                });
            })
        };

        plugin.getCard = function () {
            $.http.prerequest(function () {
            }).get("/v1/cart/card").done(function (result) {
                console.log(result);
                if (result.Status === "Success") {
                    $("input[name='CreditNumber']").val(result.Data.cardNumber)
                        .append("<input type='hidden' name='cardId' value='"+result.Data.cardId+"'>");
                    $("input[name='CreditExpiration']").val(result.Data.expiryDate);

                } else {
                    alerts("", result.Message);
                }
            }).fail(function (result) {
                console.log(result);
            });
        };

        plugin.popup = function(phone, login) {
            let $dialog = $("#specModal");

            $dialog.find("#code").val('');
            $dialog.find("#code").removeClass("error");
            $dialog.find("#code").next(".errorAlert").remove();

            $dialog.find("#code").on("keyup", function (event) {
                event.preventDefault();
                let element = $(this);
                element.removeClass(".error");
                element.next(".errorAlert").remove();
                if (element.val().length === 6 ) {
                    let data = {Phone:phone, Code:element.val()};
                    $.http.prerequest(function () {
                    }).put("/v1/cart/otp/verify", JSON.stringify(data)).done(function (result) {
                        if (result.Status === "Success") {
                            checker = true;
                            if (login) {
                                plugin.getCard();
                            }
                            $dialog.hide();
                        } else {
                            element.next(".errorAlert").remove();
                            element.addClass("error").after("<div class='errorAlert'>"+result.Message+"</div>");
                        }
                    }).fail(function (result) {
                        console.log(result);
                    });
                }
            });

            $dialog.find("#phone").text(phone.replace(/(\d{4})(\d{3})(\d{3})/, '$1-$2-$3'));
            let sec = $dialog.find('#second');
            let time = 60;
            let timer;

            function strat(phone, login) {
                timer = setInterval(function(){
                    time -= 1
                    sec.html("<em>" + time + "</em> 秒後重新寄送驗證碼");
                    if (time <= 0) {
                        sec.removeClass("m1");
                        sec.html("<a id='reCode' href='#'>重新寄送驗證碼</a>");
                        stop();
                        $("#reCode").on("click", function () {
                            plugin.submitPhone(phone, login);
                        });
                    }
                }, 1000);
            }
            function stop() {
                clearInterval(timer);
            }
            strat(phone);
            $dialog.show();
        }

        plugin.ButtonCheck = function() {
            let buttonChecker = true;
            $("form .validate").each(function (index, elem) {
                if ($(elem).val() === "" && $(elem).attr("disabled") === undefined) {
                    buttonChecker = false;
                }
                // console.log($(elem).attr("name"))
            });
            if (buttonChecker === true) {
                $("#submit").removeClass("disable").addClass("prime").attr('disabled', false);
            } else {
                $("#submit").removeClass("prime").addClass("disable").attr('disabled', true);
            }
        }


        plugin.formCheck = function() {
            $("input[name='BuyerPhone']").on("change", function (e) {
                e.preventDefault();
                let elem = $(this);
                let phone = elem.val();
                elem.removeClass("error").next(".errorAlert").remove();
                if (phone === "") {
                    elem.addClass("error").after("<div class='errorAlert'>請輸入寫手機號碼。</div>");
                    return;
                }
                let regex = /^[09]{2}[0-9]{8}$/;
                if(!regex.test(phone)) {
                    elem.addClass("error").after("<div class='errorAlert'>手機號碼輸入有誤，請確認後再重新輸入。</div>");
                    return;
                }
                plugin.submitPhone(phone, true);
            });

            $("input[name='BuyerName']").on("change", function (e){
                e.preventDefault();
                $(this).removeClass("error").next(".errorAlert").remove();
                if (($(this).val().replace(/\w/g,"")).length < 2) {
                    $(this).addClass("error").after("<div class='errorAlert'>請輸入正確姓名。</div>");
                }
            })


            $("input[name='CreditNumber']").on("keyup", function () {
                $(this).removeClass("error").next(".errorAlert").remove();
                let val = $(this).val();
                val = val.replace(/\s*/g,"");
                val = val.replace(/(\d{4})(?=\d)/g, "$1 ");
                $(this).val(val);
            })

            $("input[name='CreditExpiration']").on("keyup", function () {
                $(this).removeClass("error").next(".errorAlert").remove();
                let val = $(this).val();
                val = val.replace(/(\d{2})(?=\d)/g, "$1/");
                $(this).val( val );
            })

            $("form .validate").blur(function () {
                plugin.ButtonCheck();
            })

        };

        plugin.init();
    }

});