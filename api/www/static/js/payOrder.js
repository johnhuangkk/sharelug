$(function() {
    $.fn.payOrder = function (elm, opt) {
        let DEFAULT_OPTIONS = {
            BaseOrderId: $("#OrderId"),
            BaseStoreName: $("#StoreName"),
            BaseOrderTime: $("#OrderTime"),
            BaseBuyerName: $("#BuyerName"),
            BaseBuyerPhone: $("#BuyerPhone"),
            BasePaymentText: $(".PaymentText"),
            BaseShippingText: $(".ShippingText"),
            BaseShipping: $("#shipping"),
            BaseReceiverName: $("#ReceiverName"),
            BaseReceiverPhone: $("#ReceiverPhone"),
            BaseReceiverAdd: $("#ReceiverAdd"),
            BaseBalanceMode: $("#balanceMode"),
            BaseTransferMode: $("#Transfer"),
            BaseBankName: $("#BankName"),
            BaseBankAccount: $("#BankAccount"),
            BaseExpireDate: $("#ExpireDate"),
            BaseAmount: $("#Amount"),
            BaseBalance: $("#Balance"),

        };
        let plugin = this, options = $.extend({}, DEFAULT_OPTIONS, opt);
        let orderId = window.location.pathname.split("/").splice(-1, 1);
        console.log(orderId);
        plugin.init = function() {
            $.http.prerequest(function () {
            }).get('/v1/cart/order/' + orderId)
                .done(function (result) {
                    plugin.view(result.Data);
                }).fail(function (result) {
                window.location.href="/404";
            });
        };

        plugin.view = function(BaseData) {
            options.BaseOrderId.text(BaseData.OrderId);
            options.BaseStoreName.text(BaseData.StoreName);
            options.BaseOrderTime.text(BaseData.OrderTime);
            options.BaseBuyerName.text(BaseData.BuyerName);
            options.BaseBuyerPhone.text(BaseData.BuyerPhone);
            options.BasePaymentText.text(BaseData.Payment.Text);
            options.BaseShippingText.text(BaseData.Shipping.Text + "出貨")

            if (BaseData.Shipping.Type === "F2F") {
                options.BaseShipping.hide();
            } else {
                options.BaseShipping.show();
                options.BaseReceiverName.text(BaseData.Shipping.ReceiverName);
                options.BaseReceiverPhone.text(BaseData.Shipping.ReceiverPhone);
                options.BaseReceiverAdd.text(BaseData.Shipping.ReceiverAddress);
            }

            if (BaseData.Payment.Type === "Transfer") {
                $("#Payment").show();
                options.BaseTransferMode.show();
                options.BaseBankName.text(BaseData.Payment.BankName);
                options.BaseBankAccount.text(BaseData.Payment.BankAccount);
                options.BaseExpireDate.text(BaseData.Payment.BankExpireDate);

            }

            if (BaseData.Payment.Type === "Balance") {
                $("#Payment").show();
                options.BaseBalanceMode.show();
                options.BaseAmount.text(BaseData.Payment.OrderAmount);
                options.BaseBalance.text(BaseData.Payment.Balance);
            }


        };

        plugin.init();
    }
});