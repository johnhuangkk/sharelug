$(function () {
    $.fn.IPOST = function (elem, opt) {
        const _opt = {
            path: "/static/data/IPOST.json"
        }

        const IPOST = this, options = $.extend({}, _opt, opt);
        const source = {
            city: {},
            address: {}
        }

        const select = [
            {id: "ReceiverCity", name:"ReceiverCity", class: "validate select", text: "請選擇縣市"},
            {id: "ReceiverArea", name:"ReceiverArea", class: "validate select", text: "請選擇鄉鎮區"},
            {id: "ReceiverAddress", name:"ReceiverAddress", class: "validate", text: "請選擇地址"}
        ]

        IPOST.genOption = () => {
            const {city, address} = source
            fetch(_opt.path)
                .then((rsp) => rsp.json())
                .then((rsp) => {
                    Object.keys(rsp).forEach((c, index) => {
                        const cCountry = rsp[c].country
                        const City = rsp[c].City
                        const citys = [];
                        const opt = `<option value="${cCountry}">${cCountry}</option>`;

                        $('#ReceiverCity').append(opt);

                        City.forEach((c, index) => {
                            citys.push({city: c.city, zip: c.zip})
                            address[c.city] = c.Address
                        })
                        city[cCountry] = citys;
                    })
                    IPOST.selectChange();
                })
        }

        IPOST.init = () => {
            Object.keys(select).forEach((s, index) => {
                const id = select[index].id;
                const _class = select[index].class;
                const name = select[index].name;
                const length = $(`#${id}`).length;

                if (length === 0 && name !== "ReceiverAddress") {
                    $(`<select id="${id}" class="${_class}" name="${name}"><option value="">${select[index].text}</option></select>`).appendTo(IPOST);
                } else if(length === 0) {
                    $(`<select id="${id}" class="${_class}" name="${name}"><option value="">${select[index].text}</option></select>`).appendTo(IPOST.parent().find('div.select'));
                }
            })

            $(`<input type="hidden" id="Zipcode" name="Zipcode"/>`).appendTo(IPOST);

            IPOST.genOption();
        }

        IPOST.defaultOpt = (str) => {
            return `<option value="">${str}</option>`;
        }

        IPOST.selectChange = () => {
            IPOST.find("select").bind("change", function () {
                const $this = $(this)
                const id = $this.attr("id")
                const val = $this.val();
                const cityDefaultOpt = IPOST.defaultOpt("請選擇鄉鎮區")
                const addressDefaultOpt = IPOST.defaultOpt("請選擇地址")
                const {city, address} = source

                switch (id) {
                    case "ReceiverCity":
                        $('#ReceiverArea').empty().append(cityDefaultOpt)
                        $('#ReceiverAddress').empty().append(addressDefaultOpt)
                        if (city[val] === undefined) return;

                        city[val].forEach((c, index) => {
                            const opt = `<option value="${c.city}" data-zip="${c.zip}">${c.city}</option>`;
                            $('#ReceiverArea').append(opt);
                        })

                        break;
                    case "ReceiverArea":
                        $('#ReceiverAddress').empty().append(addressDefaultOpt)


                        if (address[val] === undefined) return;

                        $('#Zipcode').val($(this).find(':selected').data('zip'));

                        address[val].forEach((c, index) => {
                            const description = `${c.address}[${c.Location}]#${c.adm_id}`
                            const opt = `<option value="${description}" data-zip="${c.adm_id}">${description}</option>`;
                            $('#ReceiverAddress').append(opt);
                        })
                        break;

                }

            })
        }
        IPOST.init();
    }
})