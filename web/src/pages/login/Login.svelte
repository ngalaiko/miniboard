<script>
    import LoginForm from '../../components/loginform/LoginForm.svelte'

    export let api;
</script>

<div class="login">
    <LoginForm api={api} />
</div>
