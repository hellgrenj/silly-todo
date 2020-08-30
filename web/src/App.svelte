<script>
  import TodoList from "./components/TodoList.svelte";
  import Lists from "./components/Lists.svelte";
</script>

<style>
  .header {
    color:#ccc
  }
  
  h1 {
    text-transform: uppercase;
    font-size: 3em;
    font-weight: 100;
  }

</style>

<div class="container">
  <div class="row header">
    <div class="col s12">
      <h1>silly todo</h1>
    </div>
  </div>
  <div class="row">
    <div class="col s12 m3">
      <Lists />
    </div>
    <div class="col m8 offset-m1 s12">
      <TodoList />
    </div>
  </div>
</div>
