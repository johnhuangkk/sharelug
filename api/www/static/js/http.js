(function ($) {
    $.extend({
        http: {
            prerequest: function (func) {
                if (func) func();
                return this;
            },
            get: function (url) {
                return $.ajax(url, { method: "GET", dataType: "json", contentType: 'application/json; charset=utf-8' });
            },
            post: function (url, data) {
                return $.ajax(url, { method: "POST", dataType: "json", contentType: 'application/json; charset=utf-8', data: data });
            },
            put: function (url, data) {
                return $.ajax(url, { method: "PUT", dataType: "json", contentType: 'application/json; charset=utf-8', data: data });
            },
            patch: function (url, data) {
                return $.ajax(url, { method: "PATCH", dataType: "json", contentType: 'application/json; charset=utf-8', data: data });
            },
            delete: function (url) {
                return $.ajax(url, { method: "DELETE", dataType: "json", contentType: 'application/json; charset=utf-8' });
            }
        }
    });
})(jQuery);