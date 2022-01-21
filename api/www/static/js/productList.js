$(function () {
    $.fn.productList = function (elm, opt) {

        let DEFAULT_OPTIONS = {
            BaseStoreName: $("#StoreName"),
            BaseProductAmt: $("#ProductAmt"),
            BaseProfilePic: $("#ProfilePic"),
            BaseProductList: $("#productList"),
            BaseCartCount: $("#CartCount"),

        };
        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);

        let StoreId = window.location.pathname.split("/").splice(-1, 1);

        let carts;
        let BaseData;

        plugin.init = function() {
            plugin.getCartCount();

            $.http.prerequest(function () {
            }).get('/v1/store/product/list/' + StoreId)
                .done(function (result) {
                    console.log(result.Data);
                    BaseData = result.Data;
                    plugin.view();
                }).fail(function (result) {
                window.location.href="/404";
            });

            $("#pay").bind("click", function () {
                let count = parseInt($("#CartCount").html());
                if (count > 0 ){
                    window.location.href="/pay";
                } else {
                    alert("你購物車內尚無商品！");
                }

            });

            $(document).on("click", "#addCart",function () {
                let productId = $(this).data("id");
                if (carts === 0) {
                    window.location.href = "/product/"+productId;
                    return;
                }

                let $dialog = $("#specModal");
                let ProductSpecList;
                let ProductSpecId;
                let spec = 0;
                let qty = 0;
                $.each(BaseData.ProductList, function (key, value) {
                    if (value.ProductId === productId) {
                        $dialog.find(".specList").html('');
                        ProductSpecList = value.ProductSpecList;
                        $.each(ProductSpecList, function (k, v) {
                            console.log(v);
                            if (value.ProductIsSpec === 1) {
                                if (v.Quantity > 0) {
                                    $dialog.find(".specList").append("<li data-spec='" + v.ProductSpecId + "' data-quan='" + v.Quantity + "'><span>" + v.Spec + "</span></li>");
                                    spec++;
                                }
                            } else {
                                $dialog.find(".specList").append("<li data-spec='" + v.ProductSpecId + "' data-quan='" + v.Quantity + "' style='display: none'><span>" + v.Spec + "</span></li>");
                            }
                            qty = v.Quantity;
                        });
                    }
                })
                if (spec === 0) {
                    $dialog.find(".content h3").html('&nbsp<span class="close">&times;</span>');
                } else {
                    $dialog.find(".content h3").html('選擇規格<span class="close">&times;</span>');
                }
                if (qty === 0) {
                    $dialog.find(".min").text(0);
                    $dialog.find(".genBtn").removeClass("prime").addClass("disable").html("商品已售完");
                } else {
                    $dialog.find(".min").text(1);
                    $dialog.find(".genBtn").removeClass("disable").addClass("prime").html("加入結帳清單");
                }
                $dialog.find(".specList li").on("click", function() {
                    $(".specList span").removeClass("picked");
                    $(this).find("span").addClass("picked");
                    $("input[name='stock']").val($(this).data("quan")).trigger('change');
                    ProductSpecId = $(this).data("spec");
                    console.log("ProductSpecId =>", ProductSpecId);
                });
                if (ProductSpecList != null) {
                    $dialog.find("[data-spec='" + ProductSpecList[0].ProductSpecId+"']").trigger('click');
                }
                $dialog.find(".close").on("click",function (){
                    $dialog.hide();
                });

                $dialog.find(".add").bind("click", function () {
                    let data = {
                        ProductSpecId:ProductSpecId,
                        Quantity:parseInt($dialog.find(".min").html()),
                        Shipping:"",
                    }
                    plugin.addCart(data);
                });
                $dialog.show();
                $(".fAmount").QuantityBtn();

            });
        };

        plugin.addCart = function(data) {
            $.http.prerequest(function () {
            }).put('/v1/cart/add', JSON.stringify(data)).done(function (result) {
                if (result.Status === "Success") {
                    window.location.href="/pay";
                } else {
                    alerts("", result.Message);
                }
            }).fail(function (result) {
                window.location.href="/404";
            });
        };

        plugin.view = function () {
            options.BaseProfilePic.attr("src", "/static/img/tmpPic.png")
            options.BaseStoreName.append(BaseData.StoreName);
            options.BaseProductAmt.append(BaseData.ProductCount);
            options.BaseCartCount.text(carts);

            if (carts === 0) {
                $('.stickyCheckout').hide();
                $('#footer').show();
            } else {
                $('.stickyCheckout').show();
                $('#footer').hide();
            }

            if (BaseData.ProductList.length === 0) {
                options.BaseProductList.removeClass("productListWrap").addClass("errorBox").append('<p><span class="iconLoading"></span>查無資料</p>');
            } else {
                let i = 0;
                $.each(BaseData.ProductList, function (k, v) {
                     options.BaseProductList.append('<div class="prod">\n' +
                         '<p class="pic"><a href="/product/'+v.ProductId+'"><img src="'+ v.ProductImageList +'" /></a></p>\n' +
                         '<h4><a href="/product/'+v.ProductId+'">'+v.ProductName+'</a></h4>\n' +
                         '<p class="pPrice"><span class="price tw">'+v.ProductPrice+'</span></p>\n' +
                         '<p class="addToCart"><a href="javascript:void(0)" id="addCart" data-id="'+v.ProductId+'" title="加入購物清單">加入購物清單</a></p>\n' +
                         '</div>');
                     i++;
                });

                if (i < 3 ) {
                    options.BaseProductList.append('<div></div><div></div>');
                }
            }

            $('meta[property="og:site_name"]').attr("content", BaseData.StoreName);
            $('meta[property="og:url"]').attr("content", "");
            $('meta[property="og:description"]').attr("content", "");
            $('meta[property="og:image"]').attr("content", "");
        };

        plugin.getCartCount = function(){
            $.http.prerequest(function () {
            }).get('/v1/cart/count')
                .done(function (result) {
                    // console.log(result.Data);
                    carts = result.Data.Count;
                }).fail(function (result) {
                window.location.href="/404";
            });
        };


        plugin.init();
    }
});