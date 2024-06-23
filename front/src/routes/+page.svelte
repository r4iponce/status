<script>
    let description = 'Status';
    import './app.css';
    import Header from '$lib/Header.svelte';
    import Status from '$lib/Status.svelte';
</script>

<svelte:head>
    <title>Status · Home</title>
    <meta name="description" content={description}/>
</svelte:head>

<Header
        heads={[
		{ href: '/', title: 'Home', current: true },
		{ href: '/about', title: 'About', current: false }
	]}
/>

<div class="m-20 text-center text-gray-950 dark:text-gray-50">
    <div class="pb-10 text-3xl font-semibold">
        <h1>{description}</h1>
    </div>
<Status/>
</div>