$.fn.extend({


    serializeObject:function() {
        let o = {};
        let a = this.serializeArray();
        $.each(a, function() {
            if (o[this.name]) {
                if (!o[this.name].push) {
                    o[this.name] = [o[this.name]];
                }
                o[this.name].push(this.value || '');
            } else {
                o[this.name] = this.value || '';
            }
        });
        return o;
    },

    setStorage:function (data) {
        let key = "cart";
        let expire = 5000000;

        let obj = {
            time:new Date().getTime(),
            value:data,
            expire:expire,
        }
        let objStr = JSON.stringify(obj);
        localStorage.setItem(key,objStr);
    },

    getStorage:function () {
        let key = "cart";
        let expire = 5000000;

        let cartObj = [];
        let name = localStorage.getItem(key);
        let nameObj = JSON.parse(name);
        if (nameObj !== null ) {
            if (new Date().getTime() - nameObj.time >= nameObj.expire) {
                localStorage.removeItem(key)
            } else {
                cartObj = JSON.parse(nameObj.value);
            }
        }
        return cartObj;
    },

    intersect:function (a, b) {
        return $.grep(a, function(i)
        {
            return $.inArray(i, b) > -1;
        });
    },

    formatfloat:function (src, pos) {
        let num = parseFloat(src).toFixed(pos);
        num = num.toString().replace(/\$|\,/g,'');
        if (isNaN(num)) num = "0";
        let sign = (num !== (num = Math.abs(num)));
        num = Math.floor(num * 100 + 0.50000000001);
        let cents = num % 100;
        num = Math.floor(num/100).toString();
        if(cents<10) cents = "0" + cents;
        for (let i = 0; i < Math.floor((num.length-(1+i))/3); i++)
            num = num.substring(0, num.length-(4*i+3))+','+num.substring(num.length-(4*i+3));
        return (((sign)?'':'-') + num);
    },

    validatePhoneNumber:function (phone) {
        if (phone === "") {
            return false;
        } else {
            if (!/^0{0,1}(13[0-9]|15[0-9]|18[0-9]|14[0-9])[0-9]{8}$/.test(phone)) {
                return false;
            }
            return true;
        }
    },

    validateTwIdCheck:function (TwId) {
      if (TwId === "") {
          return false
      } else {
          if (!/^[A-Z]{1}[1-2]{1}[0-9]{8}$/.test(TwId)) {
              return false;
          }
          return true;
      }
    },

    validateSymbolCheck:function (input) {
        if (/[\s><,._\。\[\]\{\}\?\/\+\=\|\'\\\":;\~\!\@\#\*\$\%\^\&`\uff00-\uffff)(]+/.test(input)){
            return false;
        }
        return true;
    },

    validateChangeCount:function (input) {
        let count = input.replace(/[^\x00-\xff]/g,"**").length
        return count;
    },
});