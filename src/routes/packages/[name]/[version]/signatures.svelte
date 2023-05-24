<script lang="ts">
    export let structs:string = `Context :: struct {
  entities: Entities,
  component_map: map[typeid]Component_List,
}`
    export let procs:string = `init_ecs :: proc() -> (ctx: Context)`
    
</script>
<section class="signatures">
    <div>
        <h3>Structures:</h3>   
        <span style="color:red">todo: list the structs, hide the bodies on a triangle</span> 
        <code>
            {@html structs}
        </code>
    </div>
    <div>
        <h3>Procedures:</h3>    
        <code>
            {@html procs}
        </code>
    </div>
</section>

<style>
    .signatures{
        margin-top:1rem;
    }
</style>