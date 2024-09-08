package usuarios

import "controle_cartao/infrastructure/cadastros/usuarios"

// IUsuario define uma interface para os metodos de acesso a camada de dados
type IUsuario interface {
	CadastrarUsuario(req *usuarios.Usuario) error
	BuscarUsuario(nome *string) (*usuarios.Usuario, error)
}
