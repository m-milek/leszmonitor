import { createFileRoute, Outlet } from '@tanstack/react-router'
import {AppShell} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {TanStackRouterDevtools} from "@tanstack/react-router-devtools";
import Header from "@/components/Header.tsx";
import {Navbar} from "@/components/Navbar.tsx";

export const Route = createFileRoute('/_app')({
  component: App,
})

function App() {
    const [opened, {toggle}] = useDisclosure(false);

    return (
      <AppShell
        header={{ height: 60 }}
        navbar={{ width: 300, breakpoint: 'sm', collapsed: { mobile: !opened } }}
      >
        <AppShell.Header>
          <Header toggle={toggle} opened={opened}/>
        </AppShell.Header>
        <AppShell.Navbar>
          <Navbar/>
        </AppShell.Navbar>
        <AppShell.Main>
          <Outlet/>
        </AppShell.Main>
        <TanStackRouterDevtools/>
      </AppShell>
    )
}
