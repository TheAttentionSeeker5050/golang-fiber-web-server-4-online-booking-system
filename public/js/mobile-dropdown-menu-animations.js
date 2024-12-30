// on document ready
Zepto(function($) {
    // do something
    const dropdownMenuButton = $('#dropdown-menu-button');
    const nav = $('nav');
    
    let dropdownMenuIsOpen = false;

    dropdownMenuButton.on('click', (event) => {
        // find the #menu-icon and #menu-close-icon inside the dropdownMenuButton
        const menuIcon = $('#menu-icon');
        const menuCloseIcon = $('#menu-close-icon');
    
        if (dropdownMenuIsOpen) {
            // hide the close icon and show the menu icon
            if (menuIcon) {
                menuIcon.removeClass('hidden');
            }
    
            if (menuCloseIcon) {
                menuCloseIcon.addClass('hidden');
            }
    
            nav.removeClass('active');
            dropdownMenuIsOpen = false;
        } else {
            // hide the menu icon and show the close icon
            if (menuIcon) {
            menuIcon.addClass('hidden');
            }

            if (menuCloseIcon) {
            menuCloseIcon.removeClass('hidden');
            }

            nav.addClass('active');
            dropdownMenuIsOpen = true;
        }
    });
});
