// script.js


const Themes = Object.freeze({
    Blue: "primary",
    Gray: "secondary",
    Green: "success",
    Red: "danger",
    Yellow: "warning",
    LightBlue: "info",
    Light: "light",
    Dark: "dark"
});

const ToastPositions = Object.freeze({
    TopStart: "top-0 start-0",
    TopCenter: "top-0 start-50 translate-middle-x",
    TopEnd: "top-0 end-0",
    MiddleStart: "top-50 start-0 translate-middle-y",
    MiddleCenter: "top-50 start-50 translate-middle",
    MiddleEnd: "top-50 end-0 translate-middle-y",
    BottomStart: "bottom-0 start-0",
    BottomCenter: "bottom-0 start-50 translate-middle-x",
    BottomEnd: "bottom-0 end-0"
});

let notify = Prompt(); 

// Prompt is used to call different type alerts and popups
function Prompt() {
    // toast displays a toast popup
    let toast = function(c) {
        const { 
            title = "",
            message = "",
            theme = Themes.Light,
            position = ToastPositions.TopEnd,
            duration = 3000,
        } = c;
        
        if (message == "") {
            return
        }

        // Create a new toast element  
        var toast = document.createElement('div');
        toast.className = "toast align-items-center text-bg-" + theme;
        toast.setAttribute('role', 'alert');
        toast.setAttribute('aria-live', 'assertive');
        toast.setAttribute('aria-atomic', 'true');

        if (title !== "") {
            toast.innerHTML = `            
                <div class="toast-header">
                    <img src="..." class="rounded me-2" alt="...">
                    <strong class="me-auto">${title}</strong>                
                    <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>
                <div class="toast-body">
                    ${message}
                </div>`;
        } else {
            // TODO: workout when to add the btn-close-white attribute to the btn-close class
            toast.innerHTML = `            
                <div class="d-flex">
                    <div class="toast-body">
                        ${message}
                    </div>
                    <button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>`;
        }
    
        // Create a new toast-container element, and append toast into it
        var toastContaier = document.createElement('div');
        toastContaier.className = "toast-container position-fixed p-3 " + position;
        toastContaier.appendChild(toast)

        // Append the toast-container to the document body
        document.body.appendChild(toastContaier);

        // Show the toast
        var bsToast = new bootstrap.Toast(toast);
        bsToast.show();

        // Hide the toast after the specified duration
        setTimeout(function() {
            bsToast.hide();
            // Remove the toast from the DOM after hiding
            setTimeout(function() {
                toast.remove();
            }, 300);
        }, duration);       
    }

    return {
        toast: toast,
    }
}
        
// // notify pops up a Notie message
// function notify(msg, msgType) {
//     notie.alert({
//     type: msgType, // ['success', 'warning', 'error', 'info', 'neutral']
//     text: msg,
//     //stay: Boolean, // optional, default = false
//     //time: Number, // optional, default = 3, minimum = 1,
//     position: "bottom" // optional, default = 'top', enum: ['top', 'bottom']
//     })       
// }

// // notifyModal pops up a Swal modal
// function notifyModal(title, html, icon, button) {
//     Swal.fire({
//     title: title,
//     html: html,
//     icon: icon,
//     confirmButtonText: button
// })
// }                 

// // Prompt is used to call different type of Swal modal popups
// function Prompt2() {
//     // toast pops up a notification modal with title that disappears after 3 seconds
//     let toast = function(c) {
//         const {
//             title = "",
//             icon = "success",
//             position = "top-end",
//         } = c;

//         const Toast = Swal.mixin({
//             toast: true,
//             title: title,
//             position: position,
//             icon: icon,
//             showConfirmButton: false,
//             timer: 3000,
//             timerProgressBar: true,
//             didOpen: (toast) => {
//             toast.onmouseenter = Swal.stopTimer;
//             toast.onmouseleave = Swal.resumeTimer;
//             }
//         });

//         Toast.fire();
//     }

//     // success pops up a success modal with check-mark, title, text, footer and OK button
//     let success = function(c) {
//         const {
//             title = "",
//             text = "",
//             footer = "",         
//         } = c;

//         Swal.fire({
//             title: title,
//             text: text,
//             icon: "success",
//             footer: footer,
//             confirmButtonColor: "#1a714a" 
//         });
//     }

//     // success pops up an error modal with x-mark, title, text, footer and OK button
//     let error = function(c) {
//         const {
//             title = "",
//             text = "",
//             footer = "",         
//         } = c;

//         Swal.fire({
//             title: title,
//             text: text,
//             icon: "error",
//             footer: footer,
//             confirmButtonColor: "#1a714a" 
//         });
//     }

//     return {
//         toast: toast,
//         success: success,
//         error: error,
//     }
// }