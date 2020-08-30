<script>
  import { postJSONData } from "../utils/requestHelper.js";
  import { currentTodoListId } from "../stores/currentList.js";
  import { lastDeletedListId } from "../stores/listDeleted.js";
  /** 
   fetches all todo lists from the api 
  */
  const fetchAllLists = async () => {
    const response = await fetch("http://localhost:8080/todolist/");
    return response.json();
  };
  /** 
   creates a new todo list via the api 
  */
  const createNew = async (e) => {
    // if user hit enter
    if (e.which === 13) {
      const name = e.target.value;
      if (name) {
        const response = await postJSONData("http://localhost:8080/todolist/", {
          name,
        });
        e.target.value = "";
        allTodoLists = fetchAllLists();
      }
    }
  };
  let allTodoLists = fetchAllLists();
  /** re-populates the allTodoLists variable when a list is deleted (subscribes to changes in lastDeletedListId store) 
   which will cause the content of this component to re-render
  */
  $: {
    $lastDeletedListId;
    allTodoLists = fetchAllLists();
  }
</script>

<style>
  li {
    cursor: pointer;
  }
  input {
    max-width: 350px;
  }
</style>

{#await allTodoLists}
  <p>waiting....</p>
{:then data}
  <input
    on:keydown={createNew}
    type="text"
    placeholder="Create a new list..."
    id="newList" />
  <ul class="collection">
    {#if data}
      {#each data as list}
        <li
          class="collection-item"
          on:click={() => {
            // set value of this store to the selected id
            currentTodoListId.set(list.id);
          }}>
          {list.name.toLowerCase()}
        </li>
      {/each}
    {/if}
  </ul>

{:catch error}
  <p>an error occured</p>
{/await}
