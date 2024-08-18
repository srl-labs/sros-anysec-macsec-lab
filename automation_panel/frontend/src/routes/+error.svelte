<script lang="ts">
	import { page } from '$app/stores';

	const custom = {
		error: {
			body: "flex h-screen p-7 font-nunito",
			container: "m-auto",
			message: "rounded min-w-[400px]",
			message_header: "text-white bg-red-400 px-5 py-3 rounded-t",
			message_body: "bg-red-50 p-5 rounded-b text-red-500"
		}
	}
</script>

<div class="{custom.error.body}">
	<div class="{custom.error.container}">
		<article class="{custom.error.message}">
			<div class="{custom.error.message_header}">
				<p>Error {$page.status}</p>
			</div>
			<div class="{custom.error.message_body}">{$page.error.message}</div>
		</article>
	</div>
</div>