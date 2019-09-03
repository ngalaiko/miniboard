<script>
    export let title

    function preventNewLine(e) {
        if (!e) {
            e = window.event;
        }
        var keyCode = e.which || e.keyCode,
            target = e.target || e.srcElement;

        if (keyCode !== 13) { // Enter
            return
        }

        if (e.preventDefault) {
            e.preventDefault();
        } else {
            e.returnValue = false;
        }
    }
</script>

<div contenteditable='true' class='label' on:keydown={preventNewLine}>{title}</div>

<style>
    .label {
        display: inline-block;
        vertical-align: text-top;
        border: 1px solid;
        border-radius: 10px;
        padding: 0px 4px;
        font-size: 0.8em;
        cursor: text;
        margin-left: 3px;
    }

    .label:hover, .label:focus {
        outline-width: 0;
    }
</style>
