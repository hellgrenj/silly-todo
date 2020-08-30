<script>
  import { onDestroy } from "svelte";
  import { currentTodoListId } from "../stores/currentList.js";
  import { lastDeletedListId } from "../stores/listDeleted.js";
  import { del, uglyPatch, postJSONData } from "../utils/requestHelper.js";

  let listId;
  let todoList = fetchList();

  /** fetches the selected list*/
  function fetchList() {
    return (async () => {
      if (!listId) return null;
      const response = await fetch(`http://localhost:8080/todolist/${listId}`);
      return response.json();
    })();
  }
  // subscribes to the currentTodoListId store (see how this is set in component List.svelte)
  const unsubscribe = currentTodoListId.subscribe((value) => {
    listId = value;
    todoList = fetchList();
  });

  /** deltes a list via the api */
  const deleteList = async (e) => {
    const id = e.target.value;
    if (id) {
      const response = await del(`http://localhost:8080/todolist/${id}`);
      console.log(response.result);
      lastDeletedListId.set(id);
    }
  };

  // subscribes to lastDeletedListId store and triggers a re-render of this component if a list is deleted
  $: {
    $lastDeletedListId;
    listId = null;
    todoList = fetchList();
  }

  /** adds an item to the current todo list */
  const addItem = async (e) => {
    // if user hit enter
    if (e.which === 13) {
      const name = e.target.value;
      if (name) {
        const response = await postJSONData(
          `http://localhost:8080/todolist/${listId}/item`,
          {
            name,
            done: false,
          }
        );
        todoList = fetchList();
        
      }
    }
  };

  /** sets the done value of an item in the todo list */
  const itemCheck = async (e) => {
    const done = e.target.checked;
    const id = e.target.id;
    const response = await uglyPatch(
      `http://localhost:8080/todolist/item/${id}/${done}`
    );
    todoList = fetchList();
  };

  /** deletes an item in the todo list*/
  const deleteItem = async (id) => {
    const response = await del(`http://localhost:8080/todolist/item/${id}`);
    todoList = fetchList();
  };

  /** calls focus() on a DOM element */
  const focusOnInit = async (el) => {
    el.focus();
  }
  onDestroy(unsubscribe);
</script>

<style>
  label {
    display: inline;
    padding: 10px;
  }
  input {
    max-width: 350px;
  }
 h3 {
   font-weight: 100;
 }
 .delete-button {
   float:right;
 }
 .delete-button:hover {
   background-color:#ccc;
 }
</style>

{#await todoList then data}
  {#if data}
    <h3>
      {data.name.toLowerCase()}
      <button
        on:click={deleteList}
        class="btn btn-small waves-light btn-flat delete-button"
        value={data.id}>
        delete list
      </button>
    </h3>

    <p>
      <input
        on:keydown={addItem}
        type="text"
        placeholder="Add item to list"
        id="addItem"
        use:focusOnInit  />
    </p>

    {#if data.items}
      <ul class="collection">
        {#each data.items as item}
          <li class="collection-item">
            <label>
              <input
                on:change={itemCheck}
                type="checkbox"
                id={item.id}
                checked={item.done}
                name="items" />
              <span>{item.name}</span>
              <button
                on:click={deleteItem(item.id)}
                class="btn btn-small btn-flat delete-button"
                value={data.id}>
                del
              </button>
            </label>
          </li>
        {/each}
      </ul>
    {/if}
  {/if}
{:catch error}
  {error}
  <p>an error occured</p>
{/await}
