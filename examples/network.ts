async function getPokemon() {
  console.log('Fetching Pikachu data...');

  const response = await fetch('https://pokeapi.co/api/v2/pokemon/pikachu');

  if (!response.ok) {
    console.log('Failed to fetch!', response.status);
    return;
  }

  const data = await response.json();

  console.log('Name:', data.name);
  console.log('Weight:', data.weight);
  console.log('Base Experience:', data.base_experience);
}

getPokemon();

async function post() {
  console.log('Sending POST request to JSONPlaceholder...');

  const bodyData = {
    title: 'Brisk Runtime',
    body: 'Executing from Brisk Engine is super fast',
    userId: 1,
  };

  const response = await fetch('https://jsonplaceholder.typicode.com/posts', {
    method: 'POST',
    headers: {
      'User-Agent': 'brisk-engine-v1',
      'Content-Type': 'application/json',
      'X-Custom-Header': 'TotallyWorks',
    },
    body: JSON.stringify(bodyData),
  });

  console.log('Status Code:', response.status);
  console.log('Response OK:', response.ok);

  const json = await response.json();
  console.log('Server echoed back ID:', json.id);
  console.log('Server echoed back Title:', json.title);
}

post();
