window.addEventListener("load", function () {
    function getData() {
        var formData = new FormData(form);
        var XHR = new XMLHttpRequest();
        var url = 'https://falcon.warrensbox.com/form'
        var obj = {}

        for (var data of formData.entries()) {
            obj[data[0]] = data[1]
        }
        var nameField = document.getElementById('nameField');
        var owner_email = nameField.value

        if (owner_email === ""){
            alert("You must provide your email to test")
            return
        }

        if (!validateEmail(owner_email)){
            alert("You must provide a valid email to test")
            return
        }
        obj["owner_email"] = owner_email

        XHR.addEventListener('load', function (event) {
            alert('Message sent!\n Now, check your email to see the imformation you sent!\n');
            document.getElementById('nameField').value  = ""
            document.getElementById('contact_name').value  = ""
            document.getElementById('contact_email').value  = ""
            document.getElementById('message_content').value  = ""
        });

        // Define what happens in case of error
        XHR.addEventListener('error', function (event) {
            //alert('Message sent!\n Now, check your email to see the imformation you sent!\n');
        });
        console.log(obj)
        XHR.open('POST', url);
        XHR.send(JSON.stringify(obj));
    }

    var form = document.getElementById("myForm");
    form.addEventListener("submit", function (event) {
        event.preventDefault();
        getData();
    });

    function validateEmail(email) {
        var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return re.test(String(email).toLowerCase());
    }
});       


                            