<script>
	import { user } from '$stores/user';
</script>

<main>
    <main>
        <h1>Manage Packages</h1>
        <table>
          <thead>
            <tr>
              <th>Package Name</th>
              <th>Version</th>
              <th>Visibility</th>
              <th>Downloads</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Package 1</td>
              <td>1.0.0</td>
              <td>Public</td>
              <td>1,234</td>
              <td>
                <button>Edit</button>
                <button>Archive</button>
              </td>
            </tr>
            <!-- More rows as needed -->
          </tbody>
        </table>
      </main>
      
</main>

<style>
    table {
        width: 100%;
        text-align: left;
    }
</style>