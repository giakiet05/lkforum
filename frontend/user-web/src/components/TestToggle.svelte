<script lang="ts">
  // Một state siêu đơn giản
  let showDetails = $state(false);

  function toggle() {
    // Thay đổi state khi click
    showDetails = !showDetails;
    // Ghi log ra console để xác nhận hàm đã chạy
    console.log("Button clicked! showDetails is now:", showDetails);
  }
</script>

<div class="test-container">
  <button on:click={toggle}> Click Me to Toggle </button>

  {#if showDetails}
    <p class="details">IT WORKS! The details are now visible.</p>
  {/if}
</div>

<style>
  .test-container {
    border: 3px solid red;
    padding: 20px;
    margin: 20px;
    background-color: white;
    text-align: center;
  }
  .details {
    margin-top: 15px;
    padding: 10px;
    background-color: #e6ffed;
    color: #006421;
    font-weight: bold;
    border-radius: 5px;
  }
  button {
    font-size: 16px;
    padding: 10px 20px;
  }
</style>
