// script.js


const Themes = {
    Blue: {
        Color: "primary",
        ButtonClose: "btn-close-white",
    },
    Gray: {
        Color: "secondary",
        ButtonClose: "btn-close-white",
    },
    Green: {
        Color: "success",
        ButtonClose: "btn-close-white",
    },
    Red: {
        Color: "danger",
        ButtonClose: "btn-close-white",
    },
    Yellow: {
        Color: "warning",
        ButtonClose: "",
    },
    LightBlue: {
        Color: "info",
        ButtonClose: "",
    },
    White: {
        Color: "light",
        ButtonClose: "",
    },
    Dark: {
        Color: "dark",
        ButtonClose: "btn-close-white",
    },
};

const ToastPositions = {
    TopStart: "top-0 start-0",
    TopCenter: "top-0 start-50 translate-middle-x",
    TopEnd: "top-0 end-0",
    MiddleStart: "top-50 start-0 translate-middle-y",
    MiddleCenter: "top-50 start-50 translate-middle",
    MiddleEnd: "top-50 end-0 translate-middle-y",
    BottomStart: "bottom-0 start-0",
    BottomCenter: "bottom-0 start-50 translate-middle-x",
    BottomEnd: "bottom-0 end-0"
};

//
let notify = Prompt(); 

// Prompt is used to call different type alerts and popups
function Prompt() {
    // toast displays a toast popup
    let toast = function(c) {
        const { 
            title = "",
            message = "",
            theme = Themes.White,
            position = ToastPositions.TopEnd,
            bsIcon = "bi-info-square",
            duration = 4000,
        } = c;
        
        if (message == "") {
            return
        }

        // Create a new toast element  
        var toast = document.createElement("div");
        toast.className = "toast align-items-center text-bg-" + theme.Color;
        toast.setAttribute("role", "alert");
        toast.setAttribute("aria-live", "assertive");
        toast.setAttribute("aria-atomic", "true");

        if (title !== "") {
            toast.innerHTML = `            
                <div class="toast-header">  
                    <i class="bi ${bsIcon}" style="font-size: 5mm">&nbsp;</i>          
                    <strong class="me-auto">${title}</strong>                
                    <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>
                <div class="toast-body">                
                    ${message}
                </div>`;
        } else {
            toast.innerHTML = `            
                <div class="d-flex">
                    <div class="toast-body">
                        <div class="row align-items-center">
                            <div class="col-2 text-center">
                                <i class="bi ${bsIcon}" style="font-size: 7mm"></i>
                            </div>
                            <div class="col-10">
                                ${message}
                            </div>
                        </div>
                    </div>
                    <button type="button" class="btn-close ${theme.ButtonClose} me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>`;
        }
    
        // Create a new toast-container element, and append toast into it
        var toastContaier = document.createElement("div");
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

    // modal displays a modal popup
    let modal = function(c) {
        const { 
            title = "",
            message = "",
            bsIcon = "bi-info-square",
            confirmButtonText = "OK",
            confirmButtonTheme = Themes.Green,      
        } = c;
        
        if (message == "") {
            return
        }

        // Create a new modal element  
        var modal = document.createElement("div");
        modal.className = "modal fade";
        modal.setAttribute("tabindex", "-1");
        modal.setAttribute("aria-labelledby", "modalLabel");
        modal.setAttribute("aria-hidden", "true");        
        modal.innerHTML = `            
            <div class="modal-dialog">
                <div class="modal-content">
                <div class="modal-header">
                    <i class="bi ${bsIcon}" style="font-size: 6mm">&nbsp;</i>                    
                    <h1 class="modal-title fs-5" id="modalLabel">${title}</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    ${message}
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-${confirmButtonTheme.Color}" data-bs-dismiss="modal">${confirmButtonText}</button>
                </div>
                </div>
            </div>`;        

        // Append the modal to the document body
        document.body.appendChild(modal);

        // Show the toast
        var bsModal = new bootstrap.Modal(modal);
        bsModal.show();
    }

    return {
        toast: toast,
        modal: modal,
    }
}