import {Burger} from "@mantine/core";

interface HeaderProps {
  toggle?: () => void
  opened?: boolean;
}

export default function Header({toggle, opened}: Readonly<HeaderProps>) {
  return (
    <header>
      <nav>
        <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm"/>
        Header
      </nav>
    </header>
  )
}
