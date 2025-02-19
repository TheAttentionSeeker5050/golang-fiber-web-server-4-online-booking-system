const crudFormCallback = (e) => {
    e.preventDefault();
    var form = e.target;
    var data = {};
    $(`.base-crud-form`).serializeArray().map(function(x){
            data[x.name] = x.value;
        }
    );

    var url = form.attributes.action.value;
    var method = form.attributes.method.value;

    // make an ajax request and handle the response by adding either success or error message
    $.ajax({
        type: method,
        url: url,
        data: JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        success: function(response) {
            // add success message and the script to redirect to organizations page
            $('#success-message').html(`
                <p>Organization deleted successfully</p>
                <p>Redirecting to organizations page...</p>
            `);
        },
        error: function(xhr, status, error){
            var message = JSON.parse(xhr.response).message || "Could not add organization";
            // add error message
            $('#form-errors').html('<ul><li>' + message + '</li></ul>');
        }
    });

    setTimeout(() => {
        window.location.href = '/organizations';
    }, 3000);
}

$(document).ready(function(){
    // this is to send the form data to the server over Zepto ajax
    $('.base-crud-form').submit(crudFormCallback);
});