Rails.application.routes.draw do
  devise_for :users, path: '', path_names: {
    sign_in: 'login',
    sign_out: 'logout',
    registration: 'signup'
  },
  controllers: {
    sessions: 'users/sessions',
    registrations: 'users/registrations'
  }

  
  devise_scope :user do
    put   'signup/:id', to: 'users/registrations#update'
    patch 'signup/:id', to: 'users/registrations#update'
    delete 'signup/:id', to: 'users/registrations#destroy'  
    delete 'signup', to: 'users/registrations#destroy' 
  end
end
